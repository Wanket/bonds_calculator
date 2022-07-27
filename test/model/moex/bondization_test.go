package moex_test

import (
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/test"
	testmoex "bonds_calculator/test/model/moex"
	testdataloader "github.com/peteole/testdata-loader"
	"testing"
)

func TestDeserializeBondization(t *testing.T) {
	bondizationData := testdataloader.GetTestFile("test/data/moex/bondization.csv")
	parsedBondization := testmoex.LoadParsedBondization()

	assert, _ := test.PrepareTest(t)

	bondization, err := moex.ParseBondization(parsedBondization.ID, bondizationData)
	assert.NoError(err)

	assert.Equal(parsedBondization, bondization)
}

func BenchmarkDeserializeBondization(b *testing.B) {
	bondizationData := testdataloader.GetTestFile("test/data/moex/bondization.csv")
	parsedBondization := testmoex.LoadParsedBondization()

	for i := 0; i < b.N; i++ {
		_, _ = moex.ParseBondization(parsedBondization.ID, bondizationData)
	}

	b.ReportAllocs()
}
