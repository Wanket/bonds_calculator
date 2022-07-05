package moex

import (
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/model/moex"
	testdataloader "github.com/peteole/testdata-loader"
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

var (
	bondsData   = testdataloader.GetTestFile("test/data/moex/bond.csv")
	parsedBonds = loadParsedBonds()
)

func TestDeserializeBond(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	bonds, err := moex.ParseBondsCp1251(bondsData)
	assert.NoError(err, "unmarshalling bonds")

	assert.Equal(parsedBonds, bonds)
}

func BenchmarkDeserializeBond(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = moex.ParseBondsCp1251(bondsData)
	}

	b.ReportAllocs()
}

func loadParsedBonds() []moex.Bond {
	return []moex.Bond{
		{
			Id: "RU000A100TL1",
			SecurityPart: moex.SecurityPart{
				ShortName:     "Кузина1P01",
				Coupon:        60.41,
				NextCoupon:    ParseDate("2022-07-02"),
				AccCoupon:     38.26,
				PrevPrice:     100.08,
				Value:         4900,
				CouponPeriod:  30,
				PriceStep:     0.01,
				CouponPercent: datastuct.NewOptional(15.000),
				EndDate:       ParseDate("2023-08-26"),
			},
			MarketDataPart: moex.MarketDataPart{
				CurrentPrice: 101.8,
			},
		},
	}
}
