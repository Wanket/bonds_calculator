//go:build wireinject

//go:generate wire

package internal

import (
	"bonds_calculator/internal/api"
	"bonds_calculator/internal/controller"
	"bonds_calculator/internal/endponit"
	"bonds_calculator/internal/service"
	"bonds_calculator/internal/util"
	"github.com/benbjohnson/clock"
	"github.com/google/wire"
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

func newBondInfoService() service.BondInfoService {
	wire.Build(service.NewBondInfoService, getSingleStaticCalculatorService, getSingleStaticStoreService)

	return service.BondInfoService{}
}

func newSearchController() *controller.SearchController {
	wire.Build(controller.NewSearchController, getSingleSearchService)

	return &controller.SearchController{}
}

func newApplication() Application {
	wire.Build(NewApplication, endponit.NewRouter, endponit.NewFiberApp, getSingleSearchController)

	return Application{}
}

// singleton objects

var (
	singleSearchService           = newSearchService()
	singleBondInfoService         = newBondInfoService()
	singleStaticCalculatorService = newStaticCalculatorService()
	singleStaticStoreService      = newStaticStoreService(util.GetGlobalConfig().MoexClientQueueSize())
	singleTimeService             = newTimerService()

	searchController = newSearchController()

	app = newApplication()
)

// functions to get singletons

func GetApp() *Application {
	return &app
}

func getSingleTimeService() service.ITimerService {
	return singleTimeService
}

func getSingleStaticStoreService() service.IStaticStoreService {
	return singleStaticStoreService
}

func getSingleStaticCalculatorService() service.IStaticCalculatorService {
	return singleStaticCalculatorService
}

func getSingleSearchService() service.ISearchService {
	return singleSearchService
}

func getSingleBondInfoService() *service.BondInfoService {
	return &singleBondInfoService
}

func getSingleSearchController() *controller.SearchController {
	return searchController
}
