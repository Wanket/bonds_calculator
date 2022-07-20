package calculator

import (
	"bonds_calculator/internal/model/calculator"
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
	t.Parallel()

	assert := asserts.New(t)

	userCalculator := calculator.NewUserCalculator(userBondization, userBuyHistory)
	percent, err := userCalculator.CalcUserPercent(calculator.Maturity)
	assert.NoError(err, "CalcUserPercent calculation error")

	accuracy := 0.01
	expected := 0.125

	assert.InDelta(expected, percent, accuracy, "CalcUserPercent calculation error")
}

func TestCalcUserPercentForOneBond(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	userCalculator := calculator.NewUserCalculator(userBondization, userBuyHistory)
	percent, err := userCalculator.CalcUserPercentForOneBond(parsedBondization.Id, calculator.Maturity)

	assert.NoError(err, "CalcUserPercentForOneBond calculation error")

	accuracy := 0.01
	expected := 0.134

	assert.InDelta(expected, percent, accuracy, "CalcUserPercentForOneBond calculation error")
}

func BenchmarkCalcUserPercent(b *testing.B) {
	userCalculator := calculator.NewUserCalculator(userBondization, userBuyHistory)

	for i := 0; i < b.N; i++ {
		_, _ = userCalculator.CalcUserPercent(calculator.Maturity)
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
