package model

import (
	"bonds_calculator/internal/model"
	"bonds_calculator/internal/model/db"
	"bonds_calculator/internal/model/moex"
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

var (
	userBondization = loadUserBondization()
	userBuyHistory  = loadUserBuyHistory()
)

func TestCalcUserPercent(t *testing.T) {
	assert := asserts.New(t)

	calculator := model.NewUserCalculator(userBondization, userBuyHistory)
	percent := calculator.CalcUserPercent(model.Maturity)

	accuracy := 0.01
	expected := 0.125

	assert.InDelta(expected, percent, accuracy)
}

func TestCalcUserPercentForOneBond(t *testing.T) {
	assert := asserts.New(t)

	calculator := model.NewUserCalculator(userBondization, userBuyHistory)
	percent, err := calculator.CalcUserPercentForOneBond(parsedBondization.Id, model.Maturity)

	assert.NoError(err)

	accuracy := 0.01
	expected := 0.134

	assert.InDelta(expected, percent, accuracy)
}

func BenchmarkCalcUserPercent(b *testing.B) {
	calculator := model.NewUserCalculator(userBondization, userBuyHistory)

	for i := 0; i < b.N; i++ {
		calculator.CalcUserPercent(model.Maturity)
	}

	b.ReportAllocs()
}

func loadUserBuyHistory() []db.BuyHistory {
	return append(
		multiplyBuyHistory,
		buyHistoryVariable,
		buyHistory,
	)
}

func loadUserBondization() []moex.Bondization {
	return []moex.Bondization{
		bondizationVariable,
		parsedBondization,
	}
}
