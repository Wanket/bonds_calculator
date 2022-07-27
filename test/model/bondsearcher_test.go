package model_test

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
			ID: "FBND",
			SecurityPart: moex.SecurityPart{
				ShortName: "First Bond",
			},
		},
		{
			ID: "SBND",
			SecurityPart: moex.SecurityPart{
				ShortName: "Second Bond",
			},
		},
		{
			ID: "NBND",
			SecurityPart: moex.SecurityPart{
				ShortName: "Not B_ond",
			},
		},
		{
			ID: "NBN",
			SecurityPart: moex.SecurityPart{
				ShortName: "Not B_ond",
			},
		},
	}

	bondsSearcher := model.NewBondSearcher(bonds)

	result := bondsSearcher.Search("bonD")
	assert.Equal(sortedBondsByID(bonds[0:2]), sortedBondsByID(result))

	result = bondsSearcher.Search("BND")
	assert.Equal(sortedBondsByID(bonds[0:3]), sortedBondsByID(result))
}

func sortedBondsByID(bonds []moex.Bond) []moex.Bond {
	copyBonds := make([]moex.Bond, len(bonds))
	copy(copyBonds, bonds)

	slices.SortFunc(copyBonds, func(left, right moex.Bond) bool {
		return left.ID < right.ID
	})

	return copyBonds
}
