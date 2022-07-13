package service

import (
	"bonds_calculator/internal/model"
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastuct"
	"bonds_calculator/internal/model/datastuct/box"
	"bonds_calculator/internal/model/moex"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
	"time"
)

type SearchService struct {
	staticCalculator IStaticCalculatorService
	staticStore      IStaticStoreService

	searcher            box.ConcurrentBox[model.BondSearcher]
	searcherUpdatedTime box.ConcurrentBox[time.Time]

	reloadSearcherGroup singleflight.Group
}

func NewSearchService(staticCalculator IStaticCalculatorService, staticStore IStaticStoreService) *SearchService {
	service := SearchService{
		staticCalculator: staticCalculator,
		staticStore:      staticStore,
	}

	service.reloadSearcher()

	return &service
}

type SearchResult struct {
	Bond moex.Bond

	MaturityIncome datastuct.Optional[float64]
	CurrentIncome  datastuct.Optional[float64]
}

func (search *SearchService) Search(query string) []SearchResult {
	if updatedTime := search.searcherUpdatedTime.SafeRead(); updatedTime.Before(search.staticStore.GetBondsChangedTime()) {
		go search.reloadSearcherGroup.Do("reloadSearcher", func() (interface{}, error) {
			search.reloadSearcher()

			return nil, nil
		})
	}

	searcher := search.searcher.SafeRead()

	foundBonds := searcher.Search(query)

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
	log.Info("Reload searcher")

	bonds, updateTime := search.staticStore.GetBondsWithUpdateTime()

	search.searcher.Set(model.NewBondSearcher(bonds))
	search.searcherUpdatedTime.Set(updateTime)

	log.Info("Searcher reloaded")
}
