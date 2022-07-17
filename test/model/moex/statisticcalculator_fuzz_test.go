package moex

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/model/db"
	gofuzz "github.com/google/gofuzz"
	asserts "github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
	"math"
	"testing"
	"time"
)

func FuzzCalcStatistic(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		assert := asserts.New(t)

		fuzzer := gofuzz.NewFromGoFuzz(data)

		var income []db.Income
		fuzzer.Fuzz(&income)

		if len(income) < 2 {
			t.Skip("Income is invalid")
		}

		if !slices.IsSortedFunc(income, func(left, right db.Income) bool {
			return left.Date.Before(right.Date)
		}) {
			t.Skip("Income is invalid")
		}

		statisticCalculator := calculator.NewStatisticCalculator(income)
		results := statisticCalculator.CalcStatistic()

		slices.IsSortedFunc(results, func(left, right datastuct.Pair[time.Time, float64]) bool {
			return left.Key.Before(right.Key)
		})

		for _, result := range results {
			assert.False(math.IsNaN(result.Value), "result is NaN")

			assert.True(!income[0].Date.Truncate(time.Hour * 24).After(result.Key))
			assert.True(!result.Key.After(income[len(income)-1].Date.Truncate(time.Hour * 24).Add(time.Hour * 24)))
		}
	})
}
