package moex_test

import (
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/test"
	testmoex "bonds_calculator/test/model/moex"
	testdataloader "github.com/peteole/testdata-loader"
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

func FuzzDeserializeBondization(f *testing.F) {
	bondizationData := testdataloader.GetTestFile("test/data/moex/bondization.csv")
	parsedBondization := testmoex.LoadParsedBondization()

	f.Add(parsedBondization.ID, bondizationData)

	f.Fuzz(func(t *testing.T, id string, data []byte) {
		bondization, err := moex.ParseBondization(id, data)

		test.CheckFailed(asserts.New(t), bondization, err)
	})
}
