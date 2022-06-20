package moex

import (
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/model/moex"
	testdataloader "github.com/peteole/testdata-loader"
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

var (
	bondizationData   = testdataloader.GetTestFile("test/data/moex/bondization.csv")
	parsedBondization = loadParsedBondization()
)

func TestDeserializeBondization(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	bondization, err := moex.ParseBondization(parsedBondization.Id, bondizationData)
	assert.NoError(err, "unmarshalling bondization")

	assert.Equal(parsedBondization, bondization)
}

func BenchmarkDeserializeBondization(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = moex.ParseBondization(parsedBondization.Id, bondizationData)
	}

	b.ReportAllocs()
}

func loadParsedBondization() moex.Bondization {
	return moex.Bondization{
		Id: "RU000A100TL1",
		Amortizations: []moex.Amortization{
			{parseDate("2021-04-08"), 340},
			{parseDate("2021-05-08"), 340},
			{parseDate("2021-06-07"), 340},
			{parseDate("2021-07-07"), 340},
			{parseDate("2021-08-06"), 340},
			{parseDate("2021-09-05"), 340},
			{parseDate("2021-10-05"), 340},
			{parseDate("2021-11-04"), 340},
			{parseDate("2021-12-04"), 340},
			{parseDate("2022-01-03"), 340},
			{parseDate("2022-02-02"), 340},
			{parseDate("2022-03-04"), 340},
			{parseDate("2022-04-03"), 340},
			{parseDate("2022-05-03"), 340},
			{parseDate("2022-06-02"), 340},
			{parseDate("2022-07-02"), 340},
			{parseDate("2022-08-01"), 340},
			{parseDate("2022-08-31"), 340},
			{parseDate("2022-09-30"), 340},
			{parseDate("2022-10-30"), 340},
			{parseDate("2022-11-29"), 340},
			{parseDate("2022-12-29"), 340},
			{parseDate("2023-01-28"), 340},
			{parseDate("2023-02-27"), 340},
			{parseDate("2023-03-29"), 340},
			{parseDate("2023-04-28"), 340},
			{parseDate("2023-05-28"), 340},
			{parseDate("2023-06-27"), 340},
			{parseDate("2023-07-27"), 340},
			{parseDate("2023-08-26"), 140},
		},
		Coupons: []moex.Coupon{
			{parseDate("2019-10-16"), datastuct.NewOptional(123.29)},
			{parseDate("2019-11-15"), datastuct.NewOptional(123.29)},
			{parseDate("2019-12-15"), datastuct.NewOptional(123.29)},
			{parseDate("2020-01-14"), datastuct.NewOptional(123.29)},
			{parseDate("2020-02-13"), datastuct.NewOptional(123.29)},
			{parseDate("2020-03-14"), datastuct.NewOptional(123.29)},
			{parseDate("2020-04-13"), datastuct.NewOptional(123.29)},
			{parseDate("2020-05-13"), datastuct.NewOptional(123.29)},
			{parseDate("2020-06-12"), datastuct.NewOptional(123.29)},
			{parseDate("2020-07-12"), datastuct.NewOptional(123.29)},
			{parseDate("2020-08-11"), datastuct.NewOptional(123.29)},
			{parseDate("2020-09-10"), datastuct.NewOptional(123.29)},
			{parseDate("2020-10-10"), datastuct.NewOptional(123.29)},
			{parseDate("2020-11-09"), datastuct.NewOptional(123.29)},
			{parseDate("2020-12-09"), datastuct.NewOptional(123.29)},
			{parseDate("2021-01-08"), datastuct.NewOptional(123.29)},
			{parseDate("2021-02-07"), datastuct.NewOptional(123.29)},
			{parseDate("2021-03-09"), datastuct.NewOptional(123.29)},
			{parseDate("2021-04-08"), datastuct.NewOptional(123.29)},
			{parseDate("2021-05-08"), datastuct.NewOptional(119.1)},
			{parseDate("2021-06-07"), datastuct.NewOptional(114.9)},
			{parseDate("2021-07-07"), datastuct.NewOptional(110.71)},
			{parseDate("2021-08-06"), datastuct.NewOptional(106.52)},
			{parseDate("2021-09-05"), datastuct.NewOptional(102.33)},
			{parseDate("2021-10-05"), datastuct.NewOptional(98.14)},
			{parseDate("2021-11-04"), datastuct.NewOptional(93.95)},
			{parseDate("2021-12-04"), datastuct.NewOptional(89.75)},
			{parseDate("2022-01-03"), datastuct.NewOptional(85.56)},
			{parseDate("2022-02-02"), datastuct.NewOptional(81.37)},
			{parseDate("2022-03-04"), datastuct.NewOptional(77.18)},
			{parseDate("2022-04-03"), datastuct.NewOptional(72.99)},
			{parseDate("2022-05-03"), datastuct.NewOptional(68.79)},
			{parseDate("2022-06-02"), datastuct.NewOptional(64.6)},
			{parseDate("2022-07-02"), datastuct.NewOptional(60.41)},
			{parseDate("2022-08-01"), datastuct.NewOptional(56.22)},
			{parseDate("2022-08-31"), datastuct.NewOptional(52.03)},
			{parseDate("2022-09-30"), datastuct.NewOptional(47.84)},
			{parseDate("2022-10-30"), datastuct.NewOptional(43.64)},
			{parseDate("2022-11-29"), datastuct.NewOptional(39.45)},
			{parseDate("2022-12-29"), datastuct.NewOptional(35.26)},
			{parseDate("2023-01-28"), datastuct.NewOptional(31.07)},
			{parseDate("2023-02-27"), datastuct.NewOptional(26.88)},
			{parseDate("2023-03-29"), datastuct.NewOptional(22.68)},
			{parseDate("2023-04-28"), datastuct.NewOptional(18.49)},
			{parseDate("2023-05-28"), datastuct.NewOptional(14.3)},
			{parseDate("2023-06-27"), datastuct.NewOptional(10.11)},
			{parseDate("2023-07-27"), datastuct.NewOptional(5.92)},
			{parseDate("2023-08-26"), datastuct.NewOptional(1.73)},
		},
	}
}
