package model

import (
	"bonds_calculator/internal/model"
	"bonds_calculator/internal/model/moex"
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

func TestCalcMaturityOnePercent(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	calculator := model.NewIncomeCalculator(&parsedBondization)
	percent := calculator.CalcPercentForOneBuyHistory(buyHistory, model.Maturity)

	accuracy := 0.01
	expected := 0.124

	assert.InDelta(expected, percent, accuracy, "CalcPercentForOneBuyHistory (Maturity) calculation error")
}

func TestCalcCurrentOnePercent(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	calculator := model.NewIncomeCalculator(&parsedBondization)
	percent := calculator.CalcPercentForOneBuyHistory(buyHistory, model.Current)

	accuracy := 0.01
	expected := 0.1415

	assert.InDelta(expected, percent, accuracy, "CalcPercentForOneBuyHistory (Current) calculation error")
}

func TestCalcMultiBuyPercent(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	calculator := model.NewIncomeCalculator(&parsedBondization)
	percent := calculator.CalcPercent(multiplyBuyHistory, model.Maturity)

	accuracy := 0.01
	expected := 0.14214

	assert.InDelta(expected, percent, accuracy, "CalcPercent (MultiBuyHistory) calculation error")
}

func TestCalcVariablePercent(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	bonds, _ := moex.ParseBondsCp1251(bondsVariableData)
	bondizations, _ := moex.ParseBondization(bonds[0].Id, bondizationVariableData)

	calculator := model.NewIncomeCalculator(&bondizations)
	percent := calculator.CalcPercentForOneBuyHistory(buyHistoryVariable, model.Maturity)

	accuracy := 0.01
	expected := 0.065

	assert.InDelta(expected, percent, accuracy, "CalcPercentForOneBuyHistory (Variable) calculation error")
}

func BenchmarkCalcPercent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculator := model.NewIncomeCalculator(&parsedBondization)
		calculator.CalcPercentForOneBuyHistory(buyHistory, model.Maturity)
	}

	b.ReportAllocs()
}
