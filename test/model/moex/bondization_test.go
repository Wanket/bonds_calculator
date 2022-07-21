package moex

import (
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/test"
	testdataloader "github.com/peteole/testdata-loader"
	"testing"
)

var (
	bondizationData   = testdataloader.GetTestFile("test/data/moex/bondization.csv")
	parsedBondization = LoadParsedBondization()
)

func TestDeserializeBondization(t *testing.T) {
	assert, _ := test.PrepareTest(t)

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
