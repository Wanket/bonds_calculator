package moex

import (
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/test"
	asserts "github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/charmap"
	"testing"
)

func FuzzDeserializeBond(f *testing.F) {
	f.Add(bondsData)

	f.Fuzz(func(t *testing.T, data []byte) {
		bonds, err := moex.ParseBondsCp1251(data)

		test.CheckFailed(asserts.New(t), bonds, err)
	})
}

func FuzzDeserializeBondUtf(f *testing.F) {
	utf, _ := charmap.Windows1251.NewDecoder().Bytes(bondsData)
	f.Add(utf)

	f.Fuzz(func(t *testing.T, data []byte) {
		bonds, err := moex.ParseBondsCp1251(data)

		test.CheckFailed(asserts.New(t), bonds, err)
	})
}
