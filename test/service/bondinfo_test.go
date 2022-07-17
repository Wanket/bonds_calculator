package service

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/service"
	mock_service "bonds_calculator/internal/service/mock"
	"errors"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	asserts "github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestBondInfoService(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	assert := asserts.New(t)

	mockController := gomock.NewController(t)

	staticCalculator := mock_service.NewMockIStaticCalculatorService(mockController)
	staticStore := mock_service.NewMockIStaticStoreService(mockController)

	useBond := moex.Bond{
		Id: "1",
	}
	staticStore.EXPECT().GetBondById("1").Return(useBond, nil)

	staticStore.EXPECT().GetBondization("1").Return(moex.Bondization{
		Id: "1",
	}, nil)

	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBond, calculator.Maturity).Return(1.0, nil)
	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBond, calculator.Current).Return(2.0, nil)

	bi := service.NewBondInfoService(staticCalculator, staticStore)

	result, err := bi.GetBondInfo("1")

	assert.NoError(err)
	assert.Equal("1", result.Bond.Id)
	assert.Equal("1", result.Bondization.Id)
	assert.Equal(datastuct.NewOptional(1.0), result.MaturityIncome)
	assert.Equal(datastuct.NewOptional(2.0), result.CurrentIncome)
}

func TestBondInfoServiceBondError(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	assert := asserts.New(t)

	mockController := gomock.NewController(t)

	staticCalculator := mock_service.NewMockIStaticCalculatorService(mockController)
	staticStore := mock_service.NewMockIStaticStoreService(mockController)

	bi := service.NewBondInfoService(staticCalculator, staticStore)

	staticStore.EXPECT().GetBondById("1").Return(moex.Bond{}, errors.New("error"))

	result, err := bi.GetBondInfo("1")

	assert.Error(err)
	assert.Equal(moex.Bond{}, result.Bond)
	assert.Equal(moex.Bondization{}, result.Bondization)
	assert.Equal(datastuct.Optional[float64]{}, result.MaturityIncome)
	assert.Equal(datastuct.Optional[float64]{}, result.CurrentIncome)
}

func TestBondInfoServiceBondizationError(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	assert := asserts.New(t)

	mockController := gomock.NewController(t)

	staticCalculator := mock_service.NewMockIStaticCalculatorService(mockController)
	staticStore := mock_service.NewMockIStaticStoreService(mockController)

	bi := service.NewBondInfoService(staticCalculator, staticStore)

	useBond := moex.Bond{
		Id: "1",
	}
	staticStore.EXPECT().GetBondById("1").Return(useBond, nil)
	staticStore.EXPECT().GetBondization("1").Return(moex.Bondization{}, errors.New("error"))

	result, err := bi.GetBondInfo("1")

	assert.Error(err)
	assert.Equal(useBond, result.Bond)
	assert.Equal(moex.Bondization{}, result.Bondization)
	assert.Equal(datastuct.Optional[float64]{}, result.MaturityIncome)
	assert.Equal(datastuct.Optional[float64]{}, result.CurrentIncome)
}

func TestBondInfoServiceNoStatisticError(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	assert := asserts.New(t)

	mockController := gomock.NewController(t)

	staticCalculator := mock_service.NewMockIStaticCalculatorService(mockController)
	staticStore := mock_service.NewMockIStaticStoreService(mockController)

	bi := service.NewBondInfoService(staticCalculator, staticStore)

	useBond := moex.Bond{
		Id: "1",
	}
	staticStore.EXPECT().GetBondById("1").Return(useBond, nil)
	staticStore.EXPECT().GetBondization("1").Return(moex.Bondization{
		Id: "1",
	}, nil)
	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBond, calculator.Maturity).Return(0.0, errors.New("error"))
	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBond, calculator.Current).Return(0.0, errors.New("error"))

	result, err := bi.GetBondInfo("1")

	assert.NoError(err)
	assert.Equal(useBond, result.Bond)
	assert.Equal(moex.Bondization{Id: "1"}, result.Bondization)
	assert.Equal(datastuct.Optional[float64]{}, result.MaturityIncome)
	assert.Equal(datastuct.Optional[float64]{}, result.CurrentIncome)
}

func TestBondInfoServiceOneStatisticError(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	assert := asserts.New(t)

	mockController := gomock.NewController(t)

	staticCalculator := mock_service.NewMockIStaticCalculatorService(mockController)
	staticStore := mock_service.NewMockIStaticStoreService(mockController)

	bi := service.NewBondInfoService(staticCalculator, staticStore)

	useBond := moex.Bond{
		Id: "1",
	}
	staticStore.EXPECT().GetBondById("1").Return(useBond, nil)
	staticStore.EXPECT().GetBondization("1").Return(moex.Bondization{
		Id: "1",
	}, nil)
	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBond, calculator.Maturity).Return(1.0, errors.New("error"))
	staticCalculator.EXPECT().CalcStaticStatisticForOneBond(useBond, calculator.Current).Return(2.0, nil)

	result, err := bi.GetBondInfo("1")

	assert.NoError(err)
	assert.Equal(useBond, result.Bond)
	assert.Equal(moex.Bondization{Id: "1"}, result.Bondization)
	assert.Equal(datastuct.Optional[float64]{}, result.MaturityIncome)
	assert.Equal(datastuct.NewOptional(2.0), result.CurrentIncome)
}
