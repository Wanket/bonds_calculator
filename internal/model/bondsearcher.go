package model

import (
	"bonds_calculator/internal/model/datastuct"
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
		index.Add(strings.ToLower(bond.ShortName))
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

	uniqueBonds := datastuct.NewSet[trigram.DocID](len(ids))
	for _, id := range ids {
		uniqueBonds.Add(id / 2)
	}

	bonds := make([]moex.Bond, 0, uniqueBonds.Size())
	for inx := range uniqueBonds.Range() {
		bonds = append(bonds, searcher.bonds[inx])
	}

	return bonds
}
