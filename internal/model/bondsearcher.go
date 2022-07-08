package model

import (
	"bonds_calculator/internal/model/moex"
	"github.com/dgryski/go-trigram"
	"strings"
)

type BondSearcher struct {
	bonds []moex.Bond

	index trigram.Index
}

func NewBondSearcher(bonds []moex.Bond) BondSearcher {
	index := trigram.NewIndex([]string{})
	for _, bond := range bonds {
		index.Add(strings.ToLower(bond.Id))
		index.Add(strings.ToLower(bond.SecurityPart.ShortName))
	}

	return BondSearcher{
		bonds: bonds,
		index: index,
	}
}

func (searcher *BondSearcher) Search(query string) []moex.Bond {
	ids := searcher.index.Query(strings.ToLower(query))
	if ids == nil || len(ids) == 0 {
		return make([]moex.Bond, 0)
	}

	uniqueBonds := make(map[trigram.DocID]moex.Bond, len(ids))
	for _, id := range ids {
		uniqueBonds[id/2] = searcher.bonds[id/2]
	}

	bonds := make([]moex.Bond, 0, len(uniqueBonds))
	for _, bond := range uniqueBonds {
		bonds = append(bonds, bond)
	}

	return bonds
}
