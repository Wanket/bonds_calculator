package model

import (
	"bonds_calculator/internal/model"
	"bonds_calculator/internal/model/moex"
	"bonds_calculator/test"
	"golang.org/x/exp/slices"
	"testing"
)

func TestBondSearcher(t *testing.T) {
	assert, _ := test.PrepareTest(t)

	bonds := []moex.Bond{
		{
			Id: "FBND",
			SecurityPart: moex.SecurityPart{
				ShortName: "First Bond",
			},
		},
		{
			Id: "SBND",
			SecurityPart: moex.SecurityPart{
				ShortName: "Second Bond",
			},
		},
		{
			Id: "NBND",
			SecurityPart: moex.SecurityPart{
				ShortName: "Not B_ond",
			},
		},
		{
			Id: "NBN",
			SecurityPart: moex.SecurityPart{
				ShortName: "Not B_ond",
			},
		},
	}

	bondsSearcher := model.NewBondSearcher(bonds)

	result := bondsSearcher.Search("bonD")
	assert.Equal(sortedBondsById(bonds[0:2]), sortedBondsById(result))

	result = bondsSearcher.Search("BND")
	assert.Equal(sortedBondsById(bonds[0:3]), sortedBondsById(result))
}

func sortedBondsById(bonds []moex.Bond) []moex.Bond {
	copyBonds := make([]moex.Bond, len(bonds))
	copy(copyBonds, bonds)

	slices.SortFunc(copyBonds, func(left, right moex.Bond) bool {
		return left.Id < right.Id
	})

	return copyBonds
}
