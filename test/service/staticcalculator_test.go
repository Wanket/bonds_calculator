package service

import (
	calculator "bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/service"
	mockservice "bonds_calculator/internal/service/mock"
	"bonds_calculator/internal/util"
	"bonds_calculator/test"
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestStaticCalculator(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	calc, staticStore, _ := prepareStaticCalculator(mockController)

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
			Date:  time.Time{}.AddDate(0, 0, 10),
			Value: datastuct.NewOptional(30.0),
		},
		{
			Date:  time.Time{}.AddDate(0, 0, 20),
			Value: datastuct.Optional[float64]{},
		},
	}

	staticStore.EXPECT().GetBondization("1").Return(moex.Bondization{
		Id:            "1",
		Amortizations: expectedAmotrizations,
		Coupons:       expectedCoupons,
	}, nil)

	result, err := calc.CalcStaticStatisticForOneBond(moex.Bond{
		Id: "1",
		MarketDataPart: moex.MarketDataPart{
			CurrentPricePercent: 10.0,
		},
		SecurityPart: moex.SecurityPart{
			Coupon: 30.0,
			Value:  30.0,
		},
	}, calculator.Maturity)

	assert.NoError(err)
	assert.NotEqual(0.0, result)
}

func TestStaticCalculatorErrors(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	calc, staticStore, mockClock := prepareStaticCalculator(mockController)

	staticStore.EXPECT().GetBondization("1").Return(moex.Bondization{}, fmt.Errorf("error"))

	result, err := calc.CalcStaticStatisticForOneBond(moex.Bond{
		Id: "1",
	}, calculator.Maturity)

	assert.Error(err)
	assert.Equal(0.0, result)

	expectedAmotrizations := []moex.Amortization{
		{
			Date:  util.GetMoexNow(mockClock).AddDate(0, 0, 10),
			Value: 10,
		},
		{
			Date:  util.GetMoexNow(mockClock).AddDate(0, 0, 20),
			Value: 20,
		},
	}

	expectedCoupons := []moex.Coupon{
		{
			Date:  util.GetMoexNow(mockClock).AddDate(0, 0, 10),
			Value: datastuct.NewOptional(30.0),
		},
		{
			Date:  util.GetMoexNow(mockClock).AddDate(0, 0, 20),
			Value: datastuct.Optional[float64]{},
		},
	}

	staticStore.EXPECT().GetBondization("2").Return(moex.Bondization{
		Id:            "2",
		Amortizations: expectedAmotrizations,
		Coupons:       expectedCoupons,
	}, nil)

	result, err = calc.CalcStaticStatisticForOneBond(moex.Bond{
		Id: "2",
		MarketDataPart: moex.MarketDataPart{
			CurrentPricePercent: 10.0,
		},
		SecurityPart: moex.SecurityPart{
			Coupon: 30.0,
			Value:  20.0,
		},
	}, calculator.Maturity)

	assert.Error(err)
	assert.Equal(0.0, result)
}

func prepareStaticCalculator(mockController *gomock.Controller) (service.IStaticCalculatorService, *mockservice.MockIStaticStoreService, *clock.Mock) {
	staticStore := mockservice.NewMockIStaticStoreService(mockController)
	mockClock := clock.NewMock()

	calc := service.NewStaticCalculatorService(staticStore, mockClock)

	return calc, staticStore, mockClock
}
