package calculator_test

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/db"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/util"
	"bonds_calculator/test"
	gofuzz "github.com/google/gofuzz"
	asserts "github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

type fuzzCalcUserPercentForOneBondTestData struct {
	Bondization moex.Bondization
	BuyHistory  []db.BuyHistory
	Setting     calculator.IncomeSetting
	EndDate     time.Time
}

type fuzzCalcUserPercentTestData struct {
	Bondization []moex.Bondization
	BuyHistory  []db.BuyHistory
	Setting     calculator.IncomeSetting
	EndDate     time.Time
}

func FuzzCalcUserPercentForOneBond(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		assert := asserts.New(t)

		fuzzer := gofuzz.NewFromGoFuzz(data).NilChance(0.0)

		var fuzzCalculatorTestData fuzzCalcUserPercentForOneBondTestData
		fuzzer.Fuzz(&fuzzCalculatorTestData)

		if fuzzCalculatorTestData.Bondization.IsValid(fuzzCalculatorTestData.EndDate) != nil {
			t.Skip("Bondization is invalid")
		}

		if test.CheckBuyHistoryValid(fuzzCalculatorTestData.BuyHistory, fuzzCalculatorTestData.EndDate) != nil {
			t.Skip("BuyHistory is invalid")
		}

		userCalculator := calculator.NewUserCalculator(
			[]moex.Bondization{fuzzCalculatorTestData.Bondization},
			fuzzCalculatorTestData.BuyHistory,
		)
		result, err := userCalculator.CalcUserPercentForOneBond(
			fuzzCalculatorTestData.Bondization.ID,
			fuzzCalculatorTestData.Setting,
		)

		assert.False(err != nil && result != 0, "got error with result != 0")

		assert.False(math.IsNaN(result))
		assert.GreaterOrEqual(result, 0.0)
	})
}

func FuzzCalcUserPercent(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		assert := asserts.New(t)

		fuzzer := gofuzz.NewFromGoFuzz(data).NilChance(0.0)

		var fuzzCalculatorTestData fuzzCalcUserPercentTestData
		fuzzer.Fuzz(&fuzzCalculatorTestData)

		if util.AnyOf(fuzzCalculatorTestData.Bondization, func(bondization moex.Bondization) bool {
			return bondization.IsValid(fuzzCalculatorTestData.EndDate) != nil
		}) {
			t.Skip("Bondization is invalid")
		}

		if test.CheckBuyHistoryValid(fuzzCalculatorTestData.BuyHistory, fuzzCalculatorTestData.EndDate) != nil {
			t.Skip("BuyHistory is invalid")
		}

		userCalculator := calculator.NewUserCalculator(fuzzCalculatorTestData.Bondization, fuzzCalculatorTestData.BuyHistory)
		result, err := userCalculator.CalcUserPercent(fuzzCalculatorTestData.Setting)

		assert.False(err != nil && result != 0, "got error with result != 0")

		assert.False(math.IsNaN(result))
		assert.GreaterOrEqual(result, 0.0)
	})
}
