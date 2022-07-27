//go:generate go run github.com/mailru/easyjson/easyjson -no_std_marshalers -lower_camel_case $GOFILE
package service

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastruct"
	"bonds_calculator/internal/model/moex"
	"fmt"
	log "github.com/sirupsen/logrus"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mock/bondinfo_gen.go . IBondInfoService
type IBondInfoService interface {
	GetBondInfo(bondID string) (BondInfoResult, error)
}

type BondInfoService struct {
	staticCalculator IStaticCalculatorService
	staticStore      IStaticStoreService
}

func NewBondInfoService(staticCalculator IStaticCalculatorService, staticStore IStaticStoreService) *BondInfoService {
	return &BondInfoService{
		staticCalculator: staticCalculator,
		staticStore:      staticStore,
	}
}

//easyjson:json
type BondInfoResult struct {
	Bond        moex.Bond
	Bondization moex.Bondization

	MaturityIncome datastruct.Optional[float64]
	CurrentIncome  datastruct.Optional[float64]
}

func (infoService *BondInfoService) GetBondInfo(bondID string) (BondInfoResult, error) {
	bond, err := infoService.staticStore.GetBondByID(bondID)
	if err != nil {
		return BondInfoResult{}, fmt.Errorf("cannot get bond info, error: %w", err)
	}

	bondization, err := infoService.staticStore.GetBondization(bondID)
	if err != nil {
		return BondInfoResult{
			Bond:           bond,
			Bondization:    moex.Bondization{},
			MaturityIncome: datastruct.Optional[float64]{},
			CurrentIncome:  datastruct.Optional[float64]{},
		}, fmt.Errorf("cannot get bond info, error: %w", err)
	}

	var maturityIncome datastruct.Optional[float64]

	if maturity, err := infoService.staticCalculator.CalcStaticStatisticForOneBond(bond, calculator.Maturity); err != nil {
		log.WithFields(log.Fields{
			"bondId":     bondID,
			log.ErrorKey: err,
		}).Errorf("BondInfoService: can't calculate static maturity income")
	} else {
		maturityIncome.Set(maturity)
	}

	var currentIncome datastruct.Optional[float64]

	if current, err := infoService.staticCalculator.CalcStaticStatisticForOneBond(bond, calculator.Current); err != nil {
		log.WithFields(log.Fields{
			"bondId":     bondID,
			log.ErrorKey: err,
		}).Errorf("BondInfoService: can't calculate static current income")
	} else {
		currentIncome.Set(current)
	}

	return BondInfoResult{
		Bond:           bond,
		Bondization:    bondization,
		MaturityIncome: maturityIncome,
		CurrentIncome:  currentIncome,
	}, nil
}
