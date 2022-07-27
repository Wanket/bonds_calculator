package service_test

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastruct"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/service"
	mockservice "bonds_calculator/internal/service/mock"
	"bonds_calculator/test"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestBondInfoService(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	staticCalculator, staticStore := prepareBondInfoServiceDependencies(mockController)

	useBond := moex.Bond{
		ID: "1",
	}
	staticStore.EXPECT().GetBondByID("1").Return(useBond, nil)

	staticStore.EXPECT().GetBondization("1").Return(moex.Bondization{
		ID: "1",
	}, nil)

	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBond, calculator.Maturity).Return(1.0, nil)
	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBond, calculator.Current).Return(2.0, nil)

	bi := service.NewBondInfoService(staticCalculator, staticStore)

	result, err := bi.GetBondInfo("1")

	assert.NoError(err)
	assert.Equal("1", result.Bond.ID)
	assert.Equal("1", result.Bondization.ID)
	assert.Equal(datastruct.NewOptional(1.0), result.MaturityIncome)
	assert.Equal(datastruct.NewOptional(2.0), result.CurrentIncome)
}

func TestBondInfoServiceBondError(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	staticCalculator, staticStore := prepareBondInfoServiceDependencies(mockController)

	bi := service.NewBondInfoService(staticCalculator, staticStore)

	staticStore.EXPECT().GetBondByID("1").Return(moex.Bond{}, test.ErrTest)

	result, err := bi.GetBondInfo("1")

	assert.Error(err)
	assert.Equal(moex.Bond{}, result.Bond)
	assert.Equal(moex.Bondization{}, result.Bondization)
	assert.Equal(datastruct.Optional[float64]{}, result.MaturityIncome)
	assert.Equal(datastruct.Optional[float64]{}, result.CurrentIncome)
}

func TestBondInfoServiceBondizationError(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	staticCalculator, staticStore := prepareBondInfoServiceDependencies(mockController)

	bondInfoService := service.NewBondInfoService(staticCalculator, staticStore)

	useBond := moex.Bond{
		ID: "1",
	}
	staticStore.EXPECT().GetBondByID("1").Return(useBond, nil)
	staticStore.EXPECT().GetBondization("1").Return(moex.Bondization{}, test.ErrTest)

	result, err := bondInfoService.GetBondInfo("1")

	assert.Error(err)
	assert.Equal(useBond, result.Bond)
	assert.Equal(moex.Bondization{}, result.Bondization)
	assert.Equal(datastruct.Optional[float64]{}, result.MaturityIncome)
	assert.Equal(datastruct.Optional[float64]{}, result.CurrentIncome)
}

func TestBondInfoServiceNoStatisticError(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	staticCalculator, staticStore := prepareBondInfoServiceDependencies(mockController)

	bondInfoService := service.NewBondInfoService(staticCalculator, staticStore)

	useBond := moex.Bond{
		ID: "1",
	}
	staticStore.EXPECT().GetBondByID("1").Return(useBond, nil)
	staticStore.EXPECT().GetBondization("1").Return(moex.Bondization{
		ID: "1",
	}, nil)
	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBond, calculator.Maturity).Return(0.0, test.ErrTest)
	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBond, calculator.Current).Return(0.0, test.ErrTest)

	result, err := bondInfoService.GetBondInfo("1")

	assert.NoError(err)
	assert.Equal(useBond, result.Bond)
	assert.Equal(moex.Bondization{ID: "1"}, result.Bondization)
	assert.Equal(datastruct.Optional[float64]{}, result.MaturityIncome)
	assert.Equal(datastruct.Optional[float64]{}, result.CurrentIncome)
}

func TestBondInfoServiceOneStatisticError(t *testing.T) {
	assert, mockController := test.PrepareTest(t)

	staticCalculator, staticStore := prepareBondInfoServiceDependencies(mockController)

	bondInfoService := service.NewBondInfoService(staticCalculator, staticStore)

	useBond := moex.Bond{
		ID: "1",
	}
	staticStore.EXPECT().GetBondByID("1").Return(useBond, nil)
	staticStore.EXPECT().GetBondization("1").Return(moex.Bondization{
		ID: "1",
	}, nil)
	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBond, calculator.Maturity).Return(1.0, test.ErrTest)
	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBond, calculator.Current).Return(2.0, nil)

	result, err := bondInfoService.GetBondInfo("1")

	assert.NoError(err)
	assert.Equal(useBond, result.Bond)
	assert.Equal(moex.Bondization{ID: "1"}, result.Bondization)
	assert.Equal(datastruct.Optional[float64]{}, result.MaturityIncome)
	assert.Equal(datastruct.NewOptional(2.0), result.CurrentIncome)
}

func prepareBondInfoServiceDependencies(
	mockController *gomock.Controller,
) (*mockservice.MockIStaticCalculatorService, *mockservice.MockIStaticStoreService) {
	staticCalculator := mockservice.NewMockIStaticCalculatorService(mockController)
	staticStore := mockservice.NewMockIStaticStoreService(mockController)

	return staticCalculator, staticStore
}
