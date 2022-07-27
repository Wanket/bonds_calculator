package service_test

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastruct"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/service"
	mockservice "bonds_calculator/internal/service/mock"
	"bonds_calculator/internal/util"
	"bonds_calculator/test"
	testmoex "bonds_calculator/test/model/moex"
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
			Value: datastruct.NewOptional(30.0),
		},
		{
			Date:  time.Time{}.AddDate(0, 0, 20),
			Value: datastruct.Optional[float64]{},
		},
	}

	staticStore.EXPECT().GetBondization("1").Return(moex.Bondization{
		ID:            "1",
		Amortizations: expectedAmotrizations,
		Coupons:       expectedCoupons,
	}, nil)

	result, err := calc.CalcStaticStatisticForOneBond(moex.Bond{
		ID: "1",
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

	calc, staticStore, timeHelper := prepareStaticCalculator(mockController)

	staticStore.EXPECT().GetBondization("1").Return(moex.Bondization{}, test.ErrTest)

	result, err := calc.CalcStaticStatisticForOneBond(moex.Bond{
		ID: "1",
	}, calculator.Maturity)

	assert.Error(err)
	assert.Equal(0.0, result)

	expectedAmotrizations := []moex.Amortization{
		{
			Date:  timeHelper.GetMoexNow().AddDate(0, 0, 10),
			Value: 10,
		},
		{
			Date:  timeHelper.GetMoexNow().AddDate(0, 0, 20),
			Value: 20,
		},
	}

	expectedCoupons := []moex.Coupon{
		{
			Date:  timeHelper.GetMoexNow().AddDate(0, 0, 10),
			Value: datastruct.NewOptional(30.0),
		},
		{
			Date:  timeHelper.GetMoexNow().AddDate(0, 0, 20),
			Value: datastruct.Optional[float64]{},
		},
	}

	staticStore.EXPECT().GetBondization("2").Return(moex.Bondization{
		ID:            "2",
		Amortizations: expectedAmotrizations,
		Coupons:       expectedCoupons,
	}, nil)

	result, err = calc.CalcStaticStatisticForOneBond(moex.Bond{
		ID: "2",
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

func BenchmarkStaticCalculator(b *testing.B) {
	mockController := gomock.NewController(b)

	calc, staticStore, _ := prepareStaticCalculator(mockController)

	bondization := testmoex.LoadParsedBondization()

	staticStore.EXPECT().GetBondization(bondization.ID).AnyTimes().Return(bondization, nil)

	staticStore.EXPECT().GetBondizationsChangedTime().AnyTimes().Return(time.Time{})

	uaeBond := moex.Bond{
		ID: bondization.ID,
		MarketDataPart: moex.MarketDataPart{
			CurrentPricePercent: 10.0,
		},
		SecurityPart: moex.SecurityPart{
			Coupon: 30.0,
			Value:  30.0,
		},
	}

	for i := 0; i < b.N; i++ {
		_, _ = calc.CalcStaticStatisticForOneBond(uaeBond, calculator.Maturity)
	}
}

func prepareStaticCalculator(
	mockController *gomock.Controller,
) (*service.StaticCalculatorService, *mockservice.MockIStaticStoreService, *util.TimeHelper) {
	staticStore := mockservice.NewMockIStaticStoreService(mockController)
	mockClock := clock.NewMock()
	timeHelper := util.NewTimeHelper(mockClock)

	calc := service.NewStaticCalculatorService(staticStore, timeHelper)

	return calc, staticStore, timeHelper
}
