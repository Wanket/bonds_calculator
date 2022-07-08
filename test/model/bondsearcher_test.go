package model

import (
	"bonds_calculator/internal/model"
	"bonds_calculator/internal/model/moex"
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

func TestBondSearcher(t *testing.T) {
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
	asserts.Equal(t, bonds[0:2], result)

	result = bondsSearcher.Search("BND")
	asserts.Equal(t, bonds[0:3], result)
}
