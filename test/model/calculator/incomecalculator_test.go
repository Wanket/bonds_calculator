package calculator

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/moex"
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

func TestCalcMaturityOnePercent(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	incomeCalculator := calculator.NewIncomeCalculator(&parsedBondization)
	percent, err := incomeCalculator.CalcPercentForOneBuyHistory(buyHistory, calculator.Maturity)
	assert.NoError(err, "CalcPercentForOneBuyHistory (Maturity) calculation error")

	accuracy := 0.01
	expected := 0.124

	assert.InDelta(expected, percent, accuracy, "CalcPercentForOneBuyHistory (Maturity) calculation error")
}

func TestCalcCurrentOnePercent(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	incomeCalculator := calculator.NewIncomeCalculator(&parsedBondization)
	percent, err := incomeCalculator.CalcPercentForOneBuyHistory(buyHistory, calculator.Current)
	assert.NoError(err, "CalcPercentForOneBuyHistory (Current) calculation error")

	accuracy := 0.01
	expected := 0.1415

	assert.InDelta(expected, percent, accuracy, "CalcPercentForOneBuyHistory (Current) calculation error")
}

func TestCalcMultiBuyPercent(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	incomeCalculator := calculator.NewIncomeCalculator(&parsedBondization)
	percent, err := incomeCalculator.CalcPercent(multiplyBuyHistory, calculator.Maturity)
	assert.NoError(err, "CalcPercent (MultiBuyHistory) calculation error")

	accuracy := 0.01
	expected := 0.14214

	assert.InDelta(expected, percent, accuracy, "CalcPercent (MultiBuyHistory) calculation error")
}

func TestCalcVariablePercent(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	bonds, _ := moex.ParseBondsCp1251(bondsVariableData)
	bondizations, _ := moex.ParseBondization(bonds[0].Id, bondizationVariableData)

	incomeCalculator := calculator.NewIncomeCalculator(&bondizations)
	percent, err := incomeCalculator.CalcPercentForOneBuyHistory(buyHistoryVariable, calculator.Maturity)
	assert.NoError(err, "CalcPercentForOneBuyHistory (Variable) calculation error")

	accuracy := 0.01
	expected := 0.065

	assert.InDelta(expected, percent, accuracy, "CalcPercentForOneBuyHistory (Variable) calculation error")
}

func BenchmarkCalcPercent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		incomeCalculator := calculator.NewIncomeCalculator(&parsedBondization)
		_, _ = incomeCalculator.CalcPercentForOneBuyHistory(buyHistory, calculator.Maturity)
	}

	b.ReportAllocs()
}
