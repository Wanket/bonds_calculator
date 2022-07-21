//go:generate easyjson $GOFILE
package service

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/model/moex"
	"fmt"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=mock/bondinfo_gen.go . IBondInfoService
type IBondInfoService interface {
	GetBondInfo(bondId string) (BondInfoResult, error)
}

type BondInfoService struct {
	staticCalculator IStaticCalculatorService
	staticStore      IStaticStoreService
}

func NewBondInfoService(staticCalculator IStaticCalculatorService, staticStore IStaticStoreService) IBondInfoService {
	return &BondInfoService{
		staticCalculator: staticCalculator,
		staticStore:      staticStore,
	}
}

//easyjson:json
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
