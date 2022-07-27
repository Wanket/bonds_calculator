package calculator_test

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/test"
	testcalculator "bonds_calculator/test/model/calculator"
	testmoex "bonds_calculator/test/model/moex"
	"github.com/peteole/testdata-loader"
	"testing"
)

func TestCalcMaturityOnePercent(t *testing.T) {
	parsedBondization := testmoex.LoadParsedBondization()
	buyHistory := testcalculator.LoadBuyHistory()

	assert, _ := test.PrepareTest(t)

	incomeCalculator := calculator.NewIncomeCalculator(&parsedBondization)
	percent, err := incomeCalculator.CalcPercentForOneBuyHistory(buyHistory, calculator.Maturity)
	assert.NoError(err, "CalcPercentForOneBuyHistory (Maturity) calculation error")

	accuracy := 0.01
	expected := 0.124

	assert.InDelta(expected, percent, accuracy, "CalcPercentForOneBuyHistory (Maturity) calculation error")
}

func TestCalcCurrentOnePercent(t *testing.T) {
	parsedBondization := testmoex.LoadParsedBondization()
	buyHistory := testcalculator.LoadBuyHistory()

	assert, _ := test.PrepareTest(t)

	incomeCalculator := calculator.NewIncomeCalculator(&parsedBondization)
	percent, err := incomeCalculator.CalcPercentForOneBuyHistory(buyHistory, calculator.Current)
	assert.NoError(err, "CalcPercentForOneBuyHistory (Current) calculation error")

	accuracy := 0.01
	expected := 0.1515

	assert.InDelta(expected, percent, accuracy, "CalcPercentForOneBuyHistory (Current) calculation error")
}

func TestCalcMultiBuyPercent(t *testing.T) {
	parsedBondization := testmoex.LoadParsedBondization()
	multiplyBuyHistory := testcalculator.LoadMultiplyBuyHistory()

	assert, _ := test.PrepareTest(t)

	incomeCalculator := calculator.NewIncomeCalculator(&parsedBondization)
	percent, err := incomeCalculator.CalcPercent(multiplyBuyHistory, calculator.Maturity)
	assert.NoError(err, "CalcPercent (MultiBuyHistory) calculation error")

	accuracy := 0.01
	expected := 0.14214

	assert.InDelta(expected, percent, accuracy, "CalcPercent (MultiBuyHistory) calculation error")
}

func TestCalcVariablePercent(t *testing.T) {
	assert, _ := test.PrepareTest(t)

	buyHistoryVariable := testcalculator.LoadBuyHistoryVariable()

	bonds, err := moex.ParseBondsCp1251(testdataloader.GetTestFile("test/data/moex/bond_variable.csv"))
	assert.NoError(err)

	bondizations, err := moex.ParseBondization(
		bonds[0].ID,
		testdataloader.GetTestFile("test/data/moex/bondization_variable.csv"),
	)
	assert.NoError(err)

	incomeCalculator := calculator.NewIncomeCalculator(&bondizations)
	percent, err := incomeCalculator.CalcPercentForOneBuyHistory(buyHistoryVariable, calculator.Maturity)
	assert.NoError(err, "CalcPercentForOneBuyHistory (Variable) calculation error")

	accuracy := 0.01
	expected := 0.065

	assert.InDelta(expected, percent, accuracy, "CalcPercentForOneBuyHistory (Variable) calculation error")
}

func BenchmarkCalcPercent(b *testing.B) {
	parsedBondization := testmoex.LoadParsedBondization()
	buyHistory := testcalculator.LoadBuyHistory()

	for i := 0; i < b.N; i++ {
		incomeCalculator := calculator.NewIncomeCalculator(&parsedBondization)
		_, _ = incomeCalculator.CalcPercentForOneBuyHistory(buyHistory, calculator.Maturity)
	}

	b.ReportAllocs()
}
