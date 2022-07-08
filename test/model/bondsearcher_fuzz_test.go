package model

import (
	"bonds_calculator/internal/model"
	"bonds_calculator/internal/model/moex"
	gofuzz "github.com/google/gofuzz"
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

type fuzzBondSearcherTestData struct {
	bonds []moex.Bond

	searchString string
}

func FuzzBondSearcher(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		assert := asserts.New(t)

		fuzzer := gofuzz.NewFromGoFuzz(data)

		var testData fuzzBondSearcherTestData
		fuzzer.Fuzz(&testData)

		bondsSearcher := model.NewBondSearcher(testData.bonds)
		result := bondsSearcher.Search(testData.searchString)

		assert.NotNil(result)

		for _, bond := range result {
			assert.Contains(testData.bonds, bond)
		}
	})
}
