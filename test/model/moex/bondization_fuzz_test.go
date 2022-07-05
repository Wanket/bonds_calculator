package moex

import (
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/test"
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

func FuzzDeserializeBondization(f *testing.F) {
	f.Add(parsedBondization.Id, bondizationData)

	f.Fuzz(func(t *testing.T, id string, data []byte) {
		bondization, err := moex.ParseBondization(id, data)

		test.CheckFailed(asserts.New(t), bondization, err)
	})
}
