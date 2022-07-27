package service_test

import (
	mockapi "bonds_calculator/internal/api/mock"
	"bonds_calculator/internal/model/datastruct"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/service"
	mockservice "bonds_calculator/internal/service/mock"
	"bonds_calculator/internal/util"
	"bonds_calculator/test"
	"github.com/benbjohnson/clock"
	"github.com/golang/mock/gomock"
	"runtime"
	"testing"
	"time"
)

func validBond() moex.Bond {
	return moex.Bond{
		ID: "1",
		SecurityPart: moex.SecurityPart{
			Value:     1,
			PriceStep: 1,
			ShortName: "Test Name",
		},
		MarketDataPart: moex.MarketDataPart{},
	}
}

func TestStaticStoreCreating(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	mockClient, timerMock, _, timeHelper := prepareStaticStoreDependencies(mockController)

	mockClient.EXPECT().GetBonds().Return([]moex.Bond{}, nil)

	timerMock.EXPECT().SubscribeEvery(time.Minute*5, gomock.Any()).Return()
	timerMock.EXPECT().SubscribeEveryStartFrom(
		util.Day,
		timeHelper.GetMoexMidnight().Add(util.Day),
		gomock.Any(),
	).Return()

	now := time.Now()

	store := service.NewStaticStoreService(mockClient, timerMock, timeHelper)

	assert.Equal([]moex.Bond{}, store.GetBonds())
	assert.True(store.GetBondsChangedTime().After(now))
	assert.True(store.GetBondizationsChangedTime().After(now))
}

func TestStaticStoreBondUpdating(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	mockClient, _, clockMock, timeHelper := prepareStaticStoreDependencies(mockController)

	timer := service.NewTimerService(clockMock)
	defer timer.Close()

	mockClient.EXPECT().GetBonds().Return([]moex.Bond{}, nil)

	store := service.NewStaticStoreService(mockClient, timer, timeHelper)

	validBond := validBond()
	expectedBonds := []moex.Bond{validBond, validBond}

	expectedBonds[0].ID = "1"
	expectedBonds[0].ShortName = "First"

	expectedBonds[1].ID = "2"
	expectedBonds[1].ShortName = "Second"

	mockClient.EXPECT().GetBonds().Return(expectedBonds, nil)
	mockClient.EXPECT().GetBondization(gomock.Any()).AnyTimes().Return(moex.Bondization{}, nil)

	runtime.Gosched()

	clockMock.Add(time.Minute * 5)

	bond, err := store.GetBondByID("1")

	assert.NoError(err)
	assert.Equal("First", bond.SecurityPart.ShortName)
	assert.Equal(expectedBonds, store.GetBonds())
}

func TestStaticStoreBondizationUpdating(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	mockClient, _, clockMock, timeHelper := prepareStaticStoreDependencies(mockController)

	timer := service.NewTimerService(clockMock)
	defer timer.Close()

	useBond := []moex.Bond{validBond()}
	useBond[0].ID = "1"
	useBond[0].EndDate = time.Time{}.AddDate(0, 0, 20)

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
			Value: datastruct.NewOptional(30.0),
		},
		{
			Date:  time.Time{}.AddDate(0, 0, 40),
			Value: datastruct.Optional[float64]{},
		},
	}

	mockClient.EXPECT().GetBondization(useBond[0].ID).Return(moex.Bondization{
		ID:            "1",
		Amortizations: expectedAmotrizations,
		Coupons:       expectedCoupons,
	}, nil)

	store := service.NewStaticStoreService(mockClient, timer, timeHelper)

	bondization, err := store.GetBondization("1")

	assert.NoError(err)

	assert.Equal(expectedAmotrizations, bondization.Amortizations)
	assert.Equal(expectedCoupons, bondization.Coupons)
	assert.Equal(useBond[0].ID, bondization.ID)
}

func TestStaticStoreErrors(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	mockClient, _, clockMock, timeHelper := prepareStaticStoreDependencies(mockController)

	timer := service.NewTimerService(clockMock)
	defer timer.Close()

	useBond := []moex.Bond{validBond()}
	useBond[0].ID = "1"

	mockClient.EXPECT().GetBonds().Return(useBond, nil)

	for i := 0; i < 6; i++ {
		mockClient.EXPECT().GetBondization("1").Return(moex.Bondization{}, test.ErrTest)
	}

	store := service.NewStaticStoreService(mockClient, timer, timeHelper)

	bondization, err := store.GetBondization("1")
	assert.Error(err)
	assert.Equal(moex.Bondization{}, bondization)
}

func prepareStaticStoreDependencies(
	mockController *gomock.Controller,
) (*mockapi.MockIMoexClient, *mockservice.MockITimerService, *clock.Mock, *util.TimeHelper) {
	mockClient := mockapi.NewMockIMoexClient(mockController)
	timerMock := mockservice.NewMockITimerService(mockController)
	clockMock := clock.NewMock()
	timeHelper := util.NewTimeHelper(clockMock)

	return mockClient, timerMock, clockMock, timeHelper
}
