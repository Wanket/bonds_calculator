//go:generate go run github.com/mailru/easyjson/easyjson $GOFILE
package service

import (
	"bonds_calculator/internal/model"
	"bonds_calculator/internal/model/calculator"
	"bonds_calculator/internal/model/datastruct"
	"bonds_calculator/internal/model/datastruct/box"
	"bonds_calculator/internal/model/moex"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
	"time"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mock/search_gen.go . ISearchService
type ISearchService interface {
	Search(query string) SearchResults
}

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

		searcher:            box.ConcurrentBox[model.BondSearcher]{},
		searcherUpdatedTime: box.ConcurrentBox[time.Time]{},

		reloadSearcherGroup: singleflight.Group{},
	}

	service.reloadSearcher()

	return &service
}

//easyjson:json
type SearchResult struct {
	Bond moex.Bond

	MaturityIncome datastruct.Optional[float64]
	CurrentIncome  datastruct.Optional[float64]
}

//easyjson:json
type SearchResults []SearchResult

func (search *SearchService) Search(query string) SearchResults {
	if updatedTime := search.searcherUpdatedTime.SafeRead(); updatedTime.Before(search.staticStore.GetBondsChangedTime()) {
		_, _, _ = search.reloadSearcherGroup.Do("reloadSearcher", func() (interface{}, error) {
			search.reloadSearcher()

			return nil, nil //nolint:nilnil
		})
	}

	searcher := search.searcher.SafeRead()

	foundBonds := searcher.Search(query)

	searchResults := make([]SearchResult, 0, len(foundBonds))

	for _, bond := range foundBonds {
		var maturityIncome datastruct.Optional[float64]

		if maturity, err := search.staticCalculator.CalcStaticStatisticForOneBond(bond, calculator.Maturity); err != nil {
			log.WithFields(log.Fields{
				"bondId":     bond.ID,
				log.ErrorKey: err,
			}).Errorf("SearchService: can't calculate static maturity income")
		} else {
			maturityIncome.Set(maturity)
		}

		var currentIncome datastruct.Optional[float64]

		if current, err := search.staticCalculator.CalcStaticStatisticForOneBond(bond, calculator.Current); err != nil {
			log.WithFields(log.Fields{
				"bondId":     bond.ID,
				log.ErrorKey: err,
			}).Errorf("SearchService: can't calculate static current income")
		} else {
			currentIncome.Set(current)
		}

		searchResults = append(searchResults, SearchResult{
			Bond: bond,

			MaturityIncome: maturityIncome,
			CurrentIncome:  currentIncome,
		})
	}

	return searchResults
}

func (search *SearchService) reloadSearcher() {
	log.Info("SearchService: reload searcher")

	bonds, updateTime := search.staticStore.GetBondsWithUpdateTime()

	search.searcher.Set(model.NewBondSearcher(bonds))
	search.searcherUpdatedTime.Set(updateTime)

	log.Info("SearchService: searcher reloaded")
}
