package moex

import (
	"bonds_calculator/internal/model/moex"
	testdataloader "github.com/peteole/testdata-loader"
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

var (
	bondizationData   = testdataloader.GetTestFile("test/data/moex/bondization.csv")
	parsedBondization = LoadParsedBondization()
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
