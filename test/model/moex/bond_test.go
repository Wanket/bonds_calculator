package moex_test

import (
	"bonds_calculator/internal/model/datastruct"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/test"
	testmoex "bonds_calculator/test/model/moex"
	testdataloader "github.com/peteole/testdata-loader"
	"testing"
)

func TestDeserializeBond(t *testing.T) {
	assert, _ := test.PrepareTest(t)

	bonds, err := moex.ParseBondsCp1251(testdataloader.GetTestFile("test/data/moex/bond.csv"))
	assert.NoError(err, "unmarshalling bonds")

	assert.Equal(loadParsedBonds(), bonds)
}

func BenchmarkDeserializeBond(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = moex.ParseBondsCp1251(testdataloader.GetTestFile("test/data/moex/bond.csv"))
	}

	b.ReportAllocs()
}

func loadParsedBonds() []moex.Bond {
	return []moex.Bond{
		{
			ID: "RU000A100TL1",
			SecurityPart: moex.SecurityPart{
				ShortName:        "Кузина1P01",
				Coupon:           60.41,
				NextCoupon:       testmoex.ParseDate("2022-07-02"),
				AccCoupon:        38.26,
				PrevPricePercent: 100.08,
				Value:            4900,
				CouponPeriod:     30,
				PriceStep:        0.01,
				CouponPercent:    datastruct.NewOptional(15.000),
				EndDate:          testmoex.ParseDate("2023-08-26"),
				Currency:         "SUR",
			},
			MarketDataPart: moex.MarketDataPart{
				CurrentPricePercent: 101.8,
			},
		},
	}
}
