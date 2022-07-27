package service

import (
	modelcalculator "bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/db"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/util"
	"fmt"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mock/staticcalculator_gen.go . IStaticCalculatorService
type IStaticCalculatorService interface {
	CalcStaticStatisticForOneBond(bond moex.Bond, setting modelcalculator.IncomeSetting) (float64, error)
}

type StaticCalculatorService struct {
	staticStoreService IStaticStoreService

	timeHelper util.ITimeHelper
}

func NewStaticCalculatorService(
	staticStoreService IStaticStoreService,
	timeHelper util.ITimeHelper,
) *StaticCalculatorService {
	return &StaticCalculatorService{
		staticStoreService: staticStoreService,

		timeHelper: timeHelper,
	}
}

func (staticCalculatorService *StaticCalculatorService) CalcStaticStatisticForOneBond(
	bond moex.Bond,
	setting modelcalculator.IncomeSetting,
) (float64, error) {
	bondization, err := staticCalculatorService.staticStoreService.GetBondization(bond.ID)
	if err != nil {
		return 0, fmt.Errorf("cannot calculate statistic for bond, error: %w", err)
	}

	calculator := modelcalculator.NewIncomeCalculator(&bondization)
	result, err := calculator.CalcPercentForOneBuyHistory(db.BuyHistory{
		BondID:       bond.ID,
		Count:        1,
		Date:         staticCalculatorService.timeHelper.GetMoexNow(),
		Price:        bond.AbsoluteCurrentPrice(),
		AccCoupon:    bond.AccCoupon,
		NominalValue: bond.Value,
	}, setting)

	if err != nil {
		return 0, fmt.Errorf("cannot calculate statistic for bond, error: %w", err)
	}

	return result, nil
}
