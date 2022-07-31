package calculator

import (
	"bonds_calculator/internal/model/db"
	"bonds_calculator/internal/model/moex"
	"errors"
	"fmt"
)

var (
	errBuyHistoriesByIDNotFound = errors.New("buy histories by bond id not found")
	errBondByIDNotFound         = errors.New("bond by id not found")
)

type UserCalculator struct {
	bonds      map[string]moex.Bondization
	buyHistory map[string][]db.BuyHistory
}

func NewUserCalculator(bonds []moex.Bondization, buyHistory []db.BuyHistory) UserCalculator {
	bondsMap := make(map[string]moex.Bondization)
	for _, bond := range bonds {
		bondsMap[bond.ID] = bond
	}

	buyHistoryMap := make(map[string][]db.BuyHistory)
	for _, history := range buyHistory {
		buyHistoryMap[history.BondID] = append(buyHistoryMap[history.BondID], history)
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
		// bondCalculator lifetime is equal to bond lifetime
		bondCalculator := NewIncomeCalculator(&bond) //nolint:gosec

		histories, exist := calculator.buyHistory[bond.ID]
		if !exist {
			return 0, fmt.Errorf("CalcUserPercent: %w, id: %s", errBuyHistoriesByIDNotFound, bond.ID)
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

func (calculator *UserCalculator) CalcUserPercentForOneBond(bondID string, setting IncomeSetting) (float64, error) {
	bond, exist := calculator.bonds[bondID]
	if !exist {
		return 0, fmt.Errorf("CalcUserPercentForOneBond: %w, id: %s", errBondByIDNotFound, bondID)
	}

	histories, exist := calculator.buyHistory[bondID]
	if !exist {
		return 0, fmt.Errorf("CalcUserPercentForOneBond: %w, id: %s", errBuyHistoriesByIDNotFound, bondID)
	}

	incomeCalculator := NewIncomeCalculator(&bond)

	var (
		buyCount   uint
		sumPercent float64
	)

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
