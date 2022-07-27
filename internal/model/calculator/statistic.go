package calculator

import (
	"bonds_calculator/internal/model/datastruct"
	"bonds_calculator/internal/model/db"
	"bonds_calculator/internal/util"
	"time"
)

type StatisticCalculator struct {
	income []db.Income
}

func NewStatisticCalculator(income []db.Income) StatisticCalculator {
	return StatisticCalculator{income: income}
}

func (calculator *StatisticCalculator) CalcStatistic() []datastruct.Pair[time.Time, float64] {
	return calculator.CalcStatisticByDate(calculator.income[0].Date, calculator.income[len(calculator.income)-1].Date)
}

func (calculator *StatisticCalculator) CalcStatisticByDate(
	startDate,
	endDate time.Time,
) []datastruct.Pair[time.Time, float64] {
	incomeIndex := 0
	for incomeIndex < len(calculator.income) && calculator.income[incomeIndex].Date.Before(startDate) {
		incomeIndex++
	}

	result := make([]datastruct.Pair[time.Time, float64], 0)
	currStat := 0.0
	prevDate := time.Time{}

	for ; incomeIndex < len(calculator.income) && !calculator.income[incomeIndex].Date.After(endDate); incomeIndex++ {
		currStat += calculator.income[incomeIndex].Value

		if currDate := calculator.income[incomeIndex].Date.Truncate(util.Day); prevDate == currDate {
			result[len(result)-1].Value = currStat
		} else {
			result = append(result, datastruct.Pair[time.Time, float64]{Key: currDate, Value: currStat})
			prevDate = currDate
		}
	}

	return result
}
