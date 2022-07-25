//go:build wireinject

//go:generate go run github.com/google/wire/cmd/wire

package internal

import (
	"bonds_calculator/internal/api"
	"bonds_calculator/internal/controller"
	"bonds_calculator/internal/endponit"
	"bonds_calculator/internal/service"
	"bonds_calculator/internal/util"
	"github.com/benbjohnson/clock"
	"github.com/google/wire"
	"github.com/hymkor/go-lazy"
)

// functions which are crating New instances

func newTimerService() service.ITimerService {
	wire.Build(service.NewTimerService, clock.New)

	return &service.TimerService{}
}

func newStaticStoreService(queueSize int) service.IStaticStoreService {
	wire.Build(service.NewStaticStoreService, getSingleTimeService, api.NewMoexClient, clock.New)

	return &service.StaticStoreService{}
}

func newStaticCalculatorService() service.IStaticCalculatorService {
	wire.Build(service.NewStaticCalculatorService, getSingleStaticStoreService, clock.New)

	return &service.StaticCalculatorService{}
}

func newSearchService() service.ISearchService {
	wire.Build(service.NewSearchService, getSingleStaticCalculatorService, getSingleStaticStoreService)

	return &service.SearchService{}
}

func newBondInfoService() service.IBondInfoService {
	wire.Build(service.NewBondInfoService, getSingleStaticCalculatorService, getSingleStaticStoreService)

	return &service.BondInfoService{}
}

func newSearchController() *controller.SearchController {
	wire.Build(controller.NewSearchController, getSingleSearchService)

	return &controller.SearchController{}
}

func newBondInfoController() *controller.BondInfoController {
	wire.Build(controller.NewBondInfoController, getSingleBondInfoService)

	return &controller.BondInfoController{}
}

func newApplication() *Application {
	wire.Build(NewApplication, endponit.NewRouter, endponit.NewFiberApp, getSingleSearchController, getSingleBondInfoController)

	return &Application{}
}

// singleton objects

var (
	searchService           = lazy.New(newSearchService)
	bondInfoService         = lazy.New(newBondInfoService)
	staticCalculatorService = lazy.New(newStaticCalculatorService)
	timeService             = lazy.New(newTimerService)
	staticStoreService      = lazy.New(func() service.IStaticStoreService {
		return newStaticStoreService(util.GetGlobalConfig().MoexClientQueueSize())
	})

	searchController   = lazy.New(newSearchController)
	bondInfoController = lazy.New(newBondInfoController)

	app = lazy.New(newApplication)
)

// functions to get singletons

func GetApp() *Application {
	return app.Value()
}

func getSingleTimeService() service.ITimerService {
	return timeService.Value()
}

func getSingleStaticStoreService() service.IStaticStoreService {
	return staticStoreService.Value()
}

func getSingleStaticCalculatorService() service.IStaticCalculatorService {
	return staticCalculatorService.Value()
}

func getSingleSearchService() service.ISearchService {
	return searchService.Value()
}

func getSingleBondInfoService() service.IBondInfoService {
	return bondInfoService.Value()
}

func getSingleSearchController() *controller.SearchController {
	return searchController.Value()
}

func getSingleBondInfoController() *controller.BondInfoController {
	return bondInfoController.Value()
}
