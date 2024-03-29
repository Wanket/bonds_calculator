package moex_test

import (
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastruct"
	"bonds_calculator/internal/model/db"
	asserts "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCalcStatisticByDate(t *testing.T) {
	t.Parallel()

	income := []db.Income{
		{Date: time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC), Value: 11},
		{Date: time.Date(2018, time.January, 2, 0, 0, 0, 0, time.UTC), Value: 22},
		{Date: time.Date(2018, time.January, 2, 0, 0, 0, 0, time.UTC), Value: 33},
		{Date: time.Date(2018, time.January, 3, 0, 0, 0, 0, time.UTC), Value: 44},
		{Date: time.Date(2018, time.January, 4, 0, 0, 0, 0, time.UTC), Value: -55},
		{Date: time.Date(2018, time.January, 5, 0, 0, 0, 0, time.UTC), Value: -66},
	}

	startDate := income[1].Date
	endDate := income[len(income)-2].Date

	expectedResult := []datastruct.Pair[time.Time, float64]{
		{Key: time.Date(2018, time.January, 2, 0, 0, 0, 0, time.UTC), Value: 55},
		{Key: time.Date(2018, time.January, 3, 0, 0, 0, 0, time.UTC), Value: 99},
		{Key: time.Date(2018, time.January, 4, 0, 0, 0, 0, time.UTC), Value: 44},
	}

	statisticCalculator := calculator.NewStatisticCalculator(income)
	result := statisticCalculator.CalcStatisticByDate(startDate, endDate)

	asserts.Equal(t, expectedResult, result)
}
