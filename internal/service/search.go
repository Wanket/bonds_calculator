package service

import (
	"bonds_calculator/internal/model"
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/model/moex"
	log "github.com/sirupsen/logrus"
)

type SearchService struct {
	staticCalculator *StaticCalculatorService
	staticStore      *StaticStoreService

	searcher model.BondSearcher
}

func NewSearchService(staticCalculator *StaticCalculatorService, staticStore *StaticStoreService) SearchService {
	service := SearchService{
		staticCalculator: staticCalculator,
		staticStore:      staticStore,
	}

	service.reloadSearcher()

	return service
}

type SearchResult struct {
	Bond moex.Bond

	MaturityIncome datastuct.Optional[float64]
	CurrentIncome  datastuct.Optional[float64]
}

func (search *SearchService) Search(query string) []SearchResult {
	foundBonds := search.searcher.Search(query)

	searchResults := make([]SearchResult, 0, len(foundBonds))
	for _, bond := range foundBonds {
		searchResult := SearchResult{
			Bond: bond,
		}

		if maturity, err := search.staticCalculator.CalcStaticStatisticForOneBond(bond, calculator.Maturity); err != nil {
			log.Errorf("Can't calculate static maturity income for bond %s: %s", bond.Id, err)
		} else {
			searchResult.MaturityIncome.Set(maturity)
		}

		if current, err := search.staticCalculator.CalcStaticStatisticForOneBond(bond, calculator.Current); err != nil {
			log.Errorf("Can't calculate static current income for bond %s: %s", bond.Id, err)
		} else {
			searchResult.CurrentIncome.Set(current)
		}

		searchResults = append(searchResults, searchResult)
	}

	return searchResults
}

func (search *SearchService) reloadSearcher() {
	bonds := search.staticStore.GetBonds()

	search.searcher = model.NewBondSearcher(bonds)
}
