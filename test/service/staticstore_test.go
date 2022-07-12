package service

import (
	mockapi "bonds_calculator/internal/api/mock"
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/service"
	mockservice "bonds_calculator/internal/service/mock"
	"bonds_calculator/internal/util"
	"errors"
	"github.com/benbjohnson/clock"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	asserts "github.com/stretchr/testify/assert"
	"io/ioutil"
	"runtime"
	"testing"
	"time"
)

func TestStaticStoreCreating(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	assert := asserts.New(t)

	mockController := gomock.NewController(t)

	mockClient := mockapi.NewMockIMoexClient(mockController)
	timerMock := mockservice.NewMockITimerService(mockController)
	clockMock := clock.NewMock()

	mockClient.EXPECT().GetBonds().Return([]moex.Bond{}, nil)

	timerMock.EXPECT().SubscribeEvery(time.Minute*5, gomock.Any()).Return()
	timerMock.EXPECT().SubscribeEveryStartFrom(time.Hour*24, util.GetMoexMidnight(clockMock), gomock.Any()).Return()

	now := time.Now()

	store := service.NewStaticStoreService(mockClient, timerMock, clockMock)

	assert.Equal([]moex.Bond{}, store.GetBonds())
	assert.True(store.GetBondsChangedTime().After(now))
	assert.True(store.GetBondizationsChangedTime().After(now))
}

func TestStaticStoreBondUpdating(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	assert := asserts.New(t)

	mockController := gomock.NewController(t)

	mockClient := mockapi.NewMockIMoexClient(mockController)
	clockMock := clock.NewMock()

	timer := service.NewTimerService(clockMock)
	defer timer.Close()

	mockClient.EXPECT().GetBonds().Return([]moex.Bond{}, nil)

	store := service.NewStaticStoreService(mockClient, timer, clockMock)

	expectedBonds := []moex.Bond{
		{
			Id: "1",
			SecurityPart: moex.SecurityPart{
				ShortName: "First",
			},
		},
		{
			Id: "2",
			SecurityPart: moex.SecurityPart{
				ShortName: "Second",
			},
		},
	}

	mockClient.EXPECT().GetBonds().Return(expectedBonds, nil)

	runtime.Gosched()

	clockMock.Add(time.Minute * 5)

	bond, err := store.GetBondById("1")

	assert.NoError(err)
	assert.Equal("First", bond.SecurityPart.ShortName)
	assert.Equal(expectedBonds, store.GetBonds())
}

func TestStaticStoreBondizationUpdating(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	assert := asserts.New(t)

	mockController := gomock.NewController(t)

	mockClient := mockapi.NewMockIMoexClient(mockController)
	clockMock := clock.NewMock()

	timer := service.NewTimerService(clockMock)
	defer timer.Close()

	useBond := []moex.Bond{
		{
			Id: "1",
		},
	}

	mockClient.EXPECT().GetBonds().Return(useBond, nil)

	expectedAmotrizations := []moex.Amortization{
		{
			Date:  time.Time{}.AddDate(0, 0, 10),
			Value: 10,
		},
		{
			Date:  time.Time{}.AddDate(0, 0, 20),
			Value: 20,
		},
	}

	expectedCoupons := []moex.Coupon{
		{
			Date:  time.Time{}.AddDate(0, 0, 30),
			Value: datastuct.NewOptional(30.0),
		},
		{
			Date:  time.Time{}.AddDate(0, 0, 40),
			Value: datastuct.Optional[float64]{},
		},
	}

	mockClient.EXPECT().GetBondization(useBond[0].Id).Return(moex.Bondization{
		Id:            "1",
		Amortizations: expectedAmotrizations,
		Coupons:       expectedCoupons,
	}, nil)

	store := service.NewStaticStoreService(mockClient, timer, clockMock)

	bondization, err := store.GetBondization("1")

	assert.NoError(err)

	assert.Equal(expectedAmotrizations, bondization.Amortizations)
	assert.Equal(expectedCoupons, bondization.Coupons)
	assert.Equal(useBond[0].Id, bondization.Id)
}

func TestStaticStoreErrors(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	assert := asserts.New(t)

	mockController := gomock.NewController(t)

	mockClient := mockapi.NewMockIMoexClient(mockController)
	clockMock := clock.NewMock()

	timer := service.NewTimerService(clockMock)
	defer timer.Close()

	mockClient.EXPECT().GetBonds().Return([]moex.Bond{
		{
			Id: "1",
		},
	}, nil)

	for i := 0; i < 6; i++ {
		mockClient.EXPECT().GetBondization("1").Return(moex.Bondization{}, errors.New("error"))
	}

	store := service.NewStaticStoreService(mockClient, timer, clockMock)

	bondization, err := store.GetBondization("1")
	assert.Error(err)
	assert.Equal(moex.Bondization{}, bondization)
}
