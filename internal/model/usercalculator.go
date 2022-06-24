package model

import (
	"bonds_calculator/internal/model/db"
	"bonds_calculator/internal/model/moex"
	"fmt"
)

type UserCalculator struct {
	bonds      map[string]moex.Bondization
	buyHistory map[string][]db.BuyHistory
}

func NewUserCalculator(bonds []moex.Bondization, buyHistory []db.BuyHistory) UserCalculator {
	bondsMap := make(map[string]moex.Bondization)
	for _, bond := range bonds {
		bondsMap[bond.Id] = bond
	}

	buyHistoryMap := make(map[string][]db.BuyHistory)
	for _, history := range buyHistory {
		buyHistoryMap[history.BondId] = append(buyHistoryMap[history.BondId], history)
	}

	return UserCalculator{
		bonds:      bondsMap,
		buyHistory: buyHistoryMap,
	}
}

func (calculator *UserCalculator) CalcUserPercent(setting IncomeSetting) float64 {
	var buyCount uint
	sumPercent := 0.0
	for _, bond := range calculator.bonds {
		bondCalculator := NewIncomeCalculator(&bond)

		for _, history := range calculator.buyHistory[bond.Id] {
			sumPercent += bondCalculator.CalcPercentForOneBuyHistory(history, setting) * float64(history.Count)

			buyCount += history.Count
		}
	}

	return sumPercent / float64(buyCount)
}

func (calculator *UserCalculator) CalcUserPercentForOneBond(bondId string, setting IncomeSetting) (float64, error) {
	bond, exist := calculator.bonds[bondId]
	if !exist {
		return 0, fmt.Errorf("bond with id %s not found", bondId)
	}

	incomeCalculator := NewIncomeCalculator(&bond)

	var buyCount uint
	sumPercent := 0.0
	for _, history := range calculator.buyHistory[bondId] {
		sumPercent += incomeCalculator.CalcPercentForOneBuyHistory(history, setting) * float64(history.Count)

		buyCount += history.Count
	}

	return sumPercent / float64(buyCount), nil
}
