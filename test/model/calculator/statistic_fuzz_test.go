package calculator_test

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastruct"
	"bonds_calculator/internal/model/db"
	gofuzz "github.com/google/gofuzz"
	asserts "github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
	"testing"
	"time"
)

type fuzzCalcStatisticByDateTestData struct {
	StartDate time.Time
	EndDate   time.Time
	Income    []db.Income
}

func FuzzCalcStatistic(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		assert := asserts.New(t)

		fuzzer := gofuzz.NewFromGoFuzz(data)

		income := make([]db.Income, 0)
		fuzzer.Fuzz(&income)

		if len(income) == 0 {
			t.Skip("income is empty")
		}

		if !slices.IsSortedFunc(income, func(left, right db.Income) bool {
			return left.Date.Before(right.Date)
		}) {
			t.Skip("income is not sorted")
		}

		statisticCalculator := calculator.NewStatisticCalculator(income)
		result := statisticCalculator.CalcStatistic()

		assert.True(slices.IsSortedFunc(result, func(left, right datastruct.Pair[time.Time, float64]) bool {
			return left.Key.Sub(right.Key).Hours()/24 < 0
		}))
	})
}

func FuzzCalcStatisticByDate(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		assert := asserts.New(t)

		fuzzer := gofuzz.NewFromGoFuzz(data)

		var testData fuzzCalcStatisticByDateTestData
		fuzzer.Fuzz(&testData)

		if len(testData.Income) == 0 {
			t.Skip("income is empty")
		}

		if !slices.IsSortedFunc(testData.Income, func(left, right db.Income) bool {
			return left.Date.Before(right.Date)
		}) {
			t.Skip("income is not sorted")
		}

		statisticCalculator := calculator.NewStatisticCalculator(testData.Income)
		result := statisticCalculator.CalcStatisticByDate(testData.StartDate, testData.EndDate)

		assert.True(slices.IsSortedFunc(result, func(left, right datastruct.Pair[time.Time, float64]) bool {
			return left.Key.Sub(right.Key).Hours()/24 < 0
		}), "result is not sorted")
	})
}
