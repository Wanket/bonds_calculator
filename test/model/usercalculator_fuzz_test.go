package model

import (
	"bonds_calculator/internal/model"
	"bonds_calculator/internal/model/db"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/internal/utils"
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
	Setting     model.IncomeSetting
	EndDate     time.Time
}

type fuzzCalcUserPercentTestData struct {
	Bondization []moex.Bondization
	BuyHistory  []db.BuyHistory
	Setting     model.IncomeSetting
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

		calculator := model.NewUserCalculator([]moex.Bondization{fuzzCalculatorTestData.Bondization}, fuzzCalculatorTestData.BuyHistory)

		result, err := calculator.CalcUserPercentForOneBond(fuzzCalculatorTestData.Bondization.Id, fuzzCalculatorTestData.Setting)

		assert.False(err != nil && result != 0, "got error with result != 0")

		assert.False(math.IsNaN(result), "result is NaN")
		assert.GreaterOrEqual(result, 0.0)
	})
}

func FuzzCalcUserPercent(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		assert := asserts.New(t)

		fuzzer := gofuzz.NewFromGoFuzz(data).NilChance(0.0)

		var fuzzCalculatorTestData fuzzCalcUserPercentTestData
		fuzzer.Fuzz(&fuzzCalculatorTestData)

		if utils.AnyOf(fuzzCalculatorTestData.Bondization, func(bondization moex.Bondization) bool {
			return bondization.IsValid(fuzzCalculatorTestData.EndDate) != nil
		}) {
			t.Skip("Bondization is invalid")
		}

		if test.CheckBuyHistoryValid(fuzzCalculatorTestData.BuyHistory, fuzzCalculatorTestData.EndDate) != nil {
			t.Skip("BuyHistory is invalid")
		}

		calculator := model.NewUserCalculator(fuzzCalculatorTestData.Bondization, fuzzCalculatorTestData.BuyHistory)

		result, err := calculator.CalcUserPercent(fuzzCalculatorTestData.Setting)

		assert.False(err != nil && result != 0, "got error with result != 0")

		assert.False(math.IsNaN(result), "result is NaN")
		assert.GreaterOrEqual(result, 0.0)
	})
}
