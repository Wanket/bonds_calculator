package service

import (
	modelcalculator "bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/db"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/util"
	"fmt"
	clock "github.com/benbjohnson/clock"
)

//go:generate mockgen -destination=mock/staticcalculator_gen.go . IStaticCalculatorService
type IStaticCalculatorService interface {
	CalcStaticStatisticForOneBond(bond moex.Bond, setting modelcalculator.IncomeSetting) (float64, error)
}

type StaticCalculatorService struct {
	staticStoreService IStaticStoreService

	clock clock.Clock
}

func NewStaticCalculatorService(staticStoreService IStaticStoreService, clock clock.Clock) IStaticCalculatorService {
	return &StaticCalculatorService{
		staticStoreService: staticStoreService,
		clock:              clock,
	}
}

func (staticCalculatorService *StaticCalculatorService) CalcStaticStatisticForOneBond(bond moex.Bond, setting modelcalculator.IncomeSetting) (float64, error) {
	bondization, err := staticCalculatorService.staticStoreService.GetBondization(bond.Id)
	if err != nil {
		return 0, fmt.Errorf("cannot calculate statistic for bond, error: %v", err)
	}

	calculator := modelcalculator.NewIncomeCalculator(&bondization)
	result, err := calculator.CalcPercentForOneBuyHistory(db.BuyHistory{
		BondId:       bond.Id,
		Count:        1,
		Date:         util.GetMoexNow(staticCalculatorService.clock),
		Price:        bond.CurrentPrice,
		AccCoupon:    bond.AccCoupon,
		NominalValue: bond.Value,
	}, setting)

	if err != nil {
		return 0, fmt.Errorf("cannot calculate statistic for bond, error: %v", err)
	}

	return result, nil
}
