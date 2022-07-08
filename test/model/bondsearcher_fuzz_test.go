package model

import (
	"bonds_calculator/internal/model"
	"bonds_calculator/internal/model/moex"
	gofuzz "github.com/google/gofuzz"
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

type fuzzBondSearcherTestData struct {
	Bonds []moex.Bond

	SearchString string
}

func FuzzBondSearcher(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		assert := asserts.New(t)

		fuzzer := gofuzz.NewFromGoFuzz(data)

		var testData fuzzBondSearcherTestData
		fuzzer.Fuzz(&testData)

		if len(testData.SearchString) < 3 {
			t.Skip("Search string is too short")
		}

		bondsSearcher := model.NewBondSearcher(testData.Bonds)
		result := bondsSearcher.Search(testData.SearchString)

		assert.NotNil(result)

		for _, bond := range result {
			assert.Contains(testData.Bonds, bond)
		}
	})
}
