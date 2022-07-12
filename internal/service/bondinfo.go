package service

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/model/moex"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type BondInfoService struct {
	staticCalculator *StaticCalculatorService
	staticStore      *StaticStoreService
}

func NewBondInfoService(staticCalculator *StaticCalculatorService, staticStore *StaticStoreService) BondInfoService {
	return BondInfoService{
		staticCalculator: staticCalculator,
		staticStore:      staticStore,
	}
}

type BondInfoResult struct {
	Bond        moex.Bond
	Bondization moex.Bondization

	MaturityIncome datastuct.Optional[float64]
	CurrentIncome  datastuct.Optional[float64]
}

func (infoService *BondInfoService) GetBondInfo(bondId string) (BondInfoResult, error) {
	bondInfoResult := BondInfoResult{}

	bond, err := infoService.staticStore.GetBondById(bondId)
	if err != nil {
		return bondInfoResult, fmt.Errorf("cannot get bond info, error: %v", err)
	}

	bondInfoResult.Bond = bond

	bondization, err := infoService.staticStore.GetBondization(bondId)
	if err != nil {
		return bondInfoResult, fmt.Errorf("cannot get bond info, error: %v", err)
	}

	bondInfoResult.Bondization = bondization

	if maturity, err := infoService.staticCalculator.CalcStaticStatisticForOneBond(bond, calculator.Maturity); err != nil {
		log.Errorf("Can't calculate static maturity income for bond %s: %s", bond.Id, err)
	} else {
		bondInfoResult.MaturityIncome.Set(maturity)
	}

	if current, err := infoService.staticCalculator.CalcStaticStatisticForOneBond(bond, calculator.Current); err != nil {
		log.Errorf("Can't calculate static current income for bond %s: %s", bond.Id, err)
	} else {
		bondInfoResult.CurrentIncome.Set(current)
	}

	return bondInfoResult, nil
}
