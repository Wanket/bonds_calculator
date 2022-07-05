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

func (calculator *UserCalculator) CalcUserPercent(setting IncomeSetting) (float64, error) {
	var buyCount uint
	var sumPercent float64
	for _, bond := range calculator.bonds {
		bondCalculator := NewIncomeCalculator(&bond)

		histories, exist := calculator.buyHistory[bond.Id]
		if !exist {
			return 0, fmt.Errorf("buy histories for bond with id %s not found", bond.Id)
		}

		for _, history := range histories {
			oneHistoryPercent, err := bondCalculator.CalcPercentForOneBuyHistory(history, setting)
			if err != nil {
				return 0, err
			}

			sumPercent += oneHistoryPercent * float64(history.Count)

			buyCount += history.Count
		}
	}

	return sumPercent / float64(buyCount), nil
}

func (calculator *UserCalculator) CalcUserPercentForOneBond(bondId string, setting IncomeSetting) (float64, error) {
	bond, exist := calculator.bonds[bondId]
	if !exist {
		return 0, fmt.Errorf("bond with id %s not found", bondId)
	}

	histories, exist := calculator.buyHistory[bondId]
	if !exist {
		return 0, fmt.Errorf("buy histories for bond with id %s not found", bondId)
	}

	incomeCalculator := NewIncomeCalculator(&bond)

	var buyCount uint
	var sumPercent float64
	for _, history := range histories {
		oneHistoryPercent, err := incomeCalculator.CalcPercentForOneBuyHistory(history, setting)
		if err != nil {
			return 0, err
		}

		sumPercent += oneHistoryPercent * float64(history.Count)

		buyCount += history.Count
	}

	return sumPercent / float64(buyCount), nil
}
