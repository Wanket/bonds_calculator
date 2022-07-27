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

func newTimerService() *service.TimerService {
	wire.Build(service.NewTimerService, clock.New)

	return &service.TimerService{}
}

func newStaticStoreService() *service.StaticStoreService {
	wire.Build(service.NewStaticStoreService, getSingleMoexClient, getSingleTimeService, getSingleTimeHelper)

	return &service.StaticStoreService{}
}

func newStaticCalculatorService() *service.StaticCalculatorService {
	wire.Build(service.NewStaticCalculatorService, getSingleStaticStoreService, getSingleTimeHelper)

	return &service.StaticCalculatorService{}
}

func newSearchService() *service.SearchService {
	wire.Build(service.NewSearchService, getSingleStaticCalculatorService, getSingleStaticStoreService)

	return &service.SearchService{}
}

func newBondInfoService() *service.BondInfoService {
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
	wire.Build(
		NewApplication,
		endponit.NewRouter,
		endponit.NewFiberApp,
		getSingleSearchController,
		getSingleBondInfoController,
		getSingleConfig,
	)

	return &Application{}
}

func newGlobalConfig() *util.GlobalConfig {
	wire.Build(util.NewGlobalConfig)

	return &util.GlobalConfig{}
}

func newTimeHelper() *util.TimeHelper {
	wire.Build(util.NewTimeHelper, clock.New)

	return &util.TimeHelper{}
}

func newMoexClient(queueSize int) *api.MoexClient {
	wire.Build(api.NewMoexClient)

	return &api.MoexClient{}
}

// singleton objects

var (
	config     = lazy.New(newGlobalConfig)
	timeHelper = lazy.New(newTimeHelper)
	moexClient = lazy.New(func() *api.MoexClient {
		return newMoexClient(getSingleConfig().MoexClientQueueSize())
	})

	searchService           = lazy.New(newSearchService)
	bondInfoService         = lazy.New(newBondInfoService)
	staticCalculatorService = lazy.New(newStaticCalculatorService)
	timeService             = lazy.New(newTimerService)
	staticStoreService      = lazy.New(newStaticStoreService)

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

func getSingleConfig() util.IGlobalConfig {
	return config.Value()
}

func getSingleTimeHelper() util.ITimeHelper {
	return timeHelper.Value()
}

func getSingleMoexClient() api.IMoexClient {
	return moexClient.Value()
}
