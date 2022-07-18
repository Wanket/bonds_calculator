package service

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/model/moex"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type BondInfoService struct {
	staticCalculator IStaticCalculatorService
	staticStore      IStaticStoreService
}

func NewBondInfoService(staticCalculator IStaticCalculatorService, staticStore IStaticStoreService) BondInfoService {
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
		log.WithFields(log.Fields{
			"bondId":     bondId,
			log.ErrorKey: err,
		}).Errorf("BondInfoService: can't calculate static maturity income")
	} else {
		bondInfoResult.MaturityIncome.Set(maturity)
	}

	if current, err := infoService.staticCalculator.CalcStaticStatisticForOneBond(bond, calculator.Current); err != nil {
		log.WithFields(log.Fields{
			"bondId":     bondId,
			log.ErrorKey: err,
		}).Errorf("BondInfoService: can't calculate static current income")
	} else {
		bondInfoResult.CurrentIncome.Set(current)
	}

	return bondInfoResult, nil
}
