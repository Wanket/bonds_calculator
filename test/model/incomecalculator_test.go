package model

import (
	"bonds_calculator/internal/model"
	"bonds_calculator/internal/model/db"
	"bonds_calculator/internal/model/moex"
	testmoex "bonds_calculator/test/model/moex"
	testdataloader "github.com/peteole/testdata-loader"
	asserts "github.com/stretchr/testify/assert"
	"math"
	"testing"
)

var (
	parsedBondization = testmoex.LoadParsedBondization()

	buyHistory         = loadBuyHistory()
	multiplyBuyHistory = loadMultiplyBuyHistory()

	bondsVariableData       = testdataloader.GetTestFile("test/data/moex/bond_variable.csv")
	bondizationVariableData = testdataloader.GetTestFile("test/data/moex/bondization_variable.csv")
	buyHistoryVariable      = loadBuyHistoryVariable()
)

func TestCalcMaturityOnePercent(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	calculator := model.NewIncomeCalculator(&parsedBondization)
	percent := calculator.CalcPercentForOneBuyHistory(buyHistory, model.Maturity)

	accuracy := 0.01
	expected := 0.124

	assert.GreaterOrEqual(accuracy, math.Abs(percent-expected), "CalcPercentForOneBuyHistory (Maturity) calculation error")
}

func TestCalcCurrentOnePercent(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	calculator := model.NewIncomeCalculator(&parsedBondization)
	percent := calculator.CalcPercentForOneBuyHistory(buyHistory, model.Current)

	accuracy := 0.01
	expected := 0.1415

	assert.GreaterOrEqual(accuracy, math.Abs(percent-expected), "CalcPercentForOneBuyHistory (Current) calculation error")
}

func TestCalcMultiBuyPercent(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	calculator := model.NewIncomeCalculator(&parsedBondization)
	percent := calculator.CalcPercent(multiplyBuyHistory, model.Maturity)

	accuracy := 0.01
	expected := 0.14214

	assert.GreaterOrEqual(accuracy, math.Abs(percent-expected), "NewIncomeCalculator calculation error")
}

func TestCalcVariablePercent(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	bonds, _ := moex.ParseBondsCp1251(bondsVariableData)
	bondizations, _ := moex.ParseBondization(bonds[0].Id, bondizationVariableData)

	calculator := model.NewIncomeCalculator(&bondizations)
	percent := calculator.CalcPercentForOneBuyHistory(buyHistoryVariable, model.Maturity)

	accuracy := 0.01
	expected := 0.065

	assert.GreaterOrEqual(accuracy, math.Abs(percent-expected), "CalcPercentForOneBuyHistory (Variable) calculation error")
}

func BenchmarkCalcPercent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculator := model.NewIncomeCalculator(&parsedBondization)
		calculator.CalcPercentForOneBuyHistory(buyHistory, model.Maturity)
	}

	b.ReportAllocs()
}

func loadBuyHistory() db.BuyHistory {
	return db.BuyHistory{
		Date:         testmoex.ParseDate("2022-06-21"),
		Price:        4998.2,
		AccCoupon:    38.26,
		NominalValue: 4900,
	}
}

func loadMultiplyBuyHistory() []db.BuyHistory {
	return []db.BuyHistory{
		{
			Date:         testmoex.ParseDate("2022-06-09"),
			Price:        4999.83,
			AccCoupon:    16.11,
			NominalValue: 4900,
			Count:        2,
		},
		{
			Date:         testmoex.ParseDate("2022-05-26"),
			Price:        5203.946,
			AccCoupon:    51.68,
			NominalValue: 5240,
			Count:        3,
		},
	}
}

func loadBuyHistoryVariable() db.BuyHistory {
	return db.BuyHistory{
		Date:         testmoex.ParseDate("2022-06-23"),
		Price:        999.53,
		AccCoupon:    20.86,
		NominalValue: 1000,
	}
}
