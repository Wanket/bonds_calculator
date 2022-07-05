package model

import (
	"bonds_calculator/internal/model"
	"bonds_calculator/internal/model/db"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/test"
	gofuzz "github.com/google/gofuzz"
	asserts "github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

type fuzzCalcPercentForOneBuyHistoryTestData struct {
	Bondization moex.Bondization
	BuyHistory  db.BuyHistory
	Setting     model.IncomeSetting
	EndDate     time.Time
}

type fuzzCalcPercentTestData struct {
	Bondization moex.Bondization
	BuyHistory  []db.BuyHistory
	Setting     model.IncomeSetting
	EndDate     time.Time
}

func FuzzCalcPercentForOneBuyHistory(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		assert := asserts.New(t)

		fuzzer := gofuzz.NewFromGoFuzz(data)

		var fuzzCalculatorTestData fuzzCalcPercentForOneBuyHistoryTestData
		fuzzer.Fuzz(&fuzzCalculatorTestData)

		if fuzzCalculatorTestData.Bondization.IsValid(fuzzCalculatorTestData.EndDate) != nil {
			t.Skip("Bondization is invalid")
		}

		if test.CheckBuyHistoryValid([]db.BuyHistory{fuzzCalculatorTestData.BuyHistory}, fuzzCalculatorTestData.EndDate) != nil {
			t.Skip("BuyHistory is invalid")
		}

		calculator := model.NewIncomeCalculator(&fuzzCalculatorTestData.Bondization)

		result, err := calculator.CalcPercentForOneBuyHistory(fuzzCalculatorTestData.BuyHistory, fuzzCalculatorTestData.Setting)

		assert.False(err != nil && result != 0, "got error with result != 0")

		assert.False(math.IsNaN(result), "result is NaN")
		assert.GreaterOrEqual(result, 0.0)

		t.Logf("result: %v", calculator)
	})
}

func FuzzCalcPercent(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		assert := asserts.New(t)

		fuzzer := gofuzz.NewFromGoFuzz(data)

		var fuzzCalculatorTestData fuzzCalcPercentTestData
		fuzzer.Fuzz(&fuzzCalculatorTestData)

		if fuzzCalculatorTestData.Bondization.IsValid(fuzzCalculatorTestData.EndDate) != nil {
			t.Skip("Bondization is invalid")
		}

		if test.CheckBuyHistoryValid(fuzzCalculatorTestData.BuyHistory, fuzzCalculatorTestData.EndDate) != nil {
			t.Skip("BuyHistory is invalid")
		}

		calculator := model.NewIncomeCalculator(&fuzzCalculatorTestData.Bondization)

		result, err := calculator.CalcPercent(fuzzCalculatorTestData.BuyHistory, fuzzCalculatorTestData.Setting)

		assert.False(err != nil && result != 0, "got error with result != 0")

		assert.False(math.IsNaN(result), "result is NaN")
		assert.GreaterOrEqual(result, 0.0)
	})
}
