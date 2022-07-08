package calculator

import (
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/model/db"
	"time"
)

type StatisticCalculator struct {
	income []db.Income
}

func NewStatisticCalculator(income []db.Income) StatisticCalculator {
	return StatisticCalculator{income: income}
}

func (calculator *StatisticCalculator) CalcStatistic() []datastuct.Pair[time.Time, float64] {
	return calculator.CalcStatisticByDate(calculator.income[0].Date, calculator.income[len(calculator.income)-1].Date)
}

func (calculator *StatisticCalculator) CalcStatisticByDate(startDate, endDate time.Time) []datastuct.Pair[time.Time, float64] {
	incomeIndex := 0
	for incomeIndex < len(calculator.income) && calculator.income[incomeIndex].Date.Before(startDate) {
		incomeIndex++
	}

	result := make([]datastuct.Pair[time.Time, float64], 0)
	currStat := 0.0
	prevDate := time.Time{}
	for ; incomeIndex < len(calculator.income) && !calculator.income[incomeIndex].Date.After(endDate); incomeIndex++ {
		currStat += calculator.income[incomeIndex].Value

		if currDate := calculator.income[incomeIndex].Date.Truncate(time.Hour * 24); prevDate == currDate {
			result[len(result)-1].Value = currStat
		} else {
			result = append(result, datastuct.Pair[time.Time, float64]{currDate, currStat})
			prevDate = currDate
		}
	}

	return result
}
