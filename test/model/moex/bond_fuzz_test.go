package moex_test

import (
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/test"
	"github.com/peteole/testdata-loader"
	asserts "github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/charmap"
	"testing"
)

func FuzzDeserializeBond(f *testing.F) {
	f.Add(testdataloader.GetTestFile("test/data/moex/bond.csv"))

	f.Fuzz(func(t *testing.T, data []byte) {
		bonds, err := moex.ParseBondsCp1251(data)

		test.CheckFailed(asserts.New(t), bonds, err)
	})
}

func FuzzDeserializeBondUtf(f *testing.F) {
	utf, err := charmap.Windows1251.NewDecoder().Bytes(testdataloader.GetTestFile("test/data/moex/bond.csv"))
	asserts.NoError(f, err)

	f.Add(utf)

	f.Fuzz(func(t *testing.T, data []byte) {
		bonds, err := moex.ParseBondsCp1251(data)

		test.CheckFailed(asserts.New(t), bonds, err)
	})
}
