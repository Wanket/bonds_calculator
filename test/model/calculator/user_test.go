package calculator_test

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/db"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/test"
	testcalculator "bonds_calculator/test/model/calculator"
	testmoex "bonds_calculator/test/model/moex"
	"testing"
)

func TestCalcUserPercent(t *testing.T) {
	assert, _ := test.PrepareTest(t)

	userCalculator := calculator.NewUserCalculator(loadUserBondization(), loadUserBuyHistory())
	percent, err := userCalculator.CalcUserPercent(calculator.Maturity)
	assert.NoError(err)

	accuracy := 0.01
	expected := 0.125

	assert.InDelta(expected, percent, accuracy)
}

func TestCalcUserPercentForOneBond(t *testing.T) {
	assert, _ := test.PrepareTest(t)

	userCalculator := calculator.NewUserCalculator(loadUserBondization(), loadUserBuyHistory())
	percent, err := userCalculator.CalcUserPercentForOneBond(testmoex.LoadParsedBondization().ID, calculator.Maturity)

	assert.NoError(err, "CalcUserPercentForOneBond calculation error")

	accuracy := 0.01
	expected := 0.134

	assert.InDelta(expected, percent, accuracy, "CalcUserPercentForOneBond calculation error")
}

func BenchmarkCalcUserPercent(b *testing.B) {
	userCalculator := calculator.NewUserCalculator(loadUserBondization(), loadUserBuyHistory())

	for i := 0; i < b.N; i++ {
		_, _ = userCalculator.CalcUserPercent(calculator.Maturity)
	}

	b.ReportAllocs()
}

func loadUserBuyHistory() []db.BuyHistory {
	return append(
		testcalculator.LoadMultiplyBuyHistory(),
		testcalculator.LoadBuyHistoryVariable(),
		testcalculator.LoadBuyHistory(),
	)
}

func loadUserBondization() []moex.Bondization {
	return []moex.Bondization{
		testcalculator.LoadBondizationVariable(),
		testmoex.LoadParsedBondization(),
	}
}
