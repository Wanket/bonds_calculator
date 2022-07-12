//go:build wireinject

//go:generate wire

package internal

import (
	"bonds_calculator/internal/api"
	"bonds_calculator/internal/service"
	"bonds_calculator/internal/util"
	"github.com/benbjohnson/clock"
	"github.com/google/wire"
)

func NewTimerService() service.ITimerService {
	wire.Build(service.NewTimerService, clock.New)

	return &service.TimerService{}
}

var singleTimeService = NewTimerService()

func newSingleTimeService() service.ITimerService {
	return singleTimeService
}

func NewStaticStoreService(queueSize int) *service.StaticStoreService {
	wire.Build(service.NewStaticStoreService, newSingleTimeService, api.NewMoexClient, clock.New)

	return &service.StaticStoreService{}
}

var singleStaticStoreService = NewStaticStoreService(util.GetGlobalConfig().MoexClientQueueSize())

func newSingleStaticStoreService() *service.StaticStoreService {
	return singleStaticStoreService
}

func NewStaticCalculatorService() service.StaticCalculatorService {
	wire.Build(service.NewStaticCalculatorService, newSingleStaticStoreService, clock.New)

	return service.StaticCalculatorService{}
}

var singleStaticCalculatorService = NewStaticCalculatorService()

func newSingleStaticCalculatorService() *service.StaticCalculatorService {
	return &singleStaticCalculatorService
}

func NewSearchService() service.SearchService {
	wire.Build(service.NewSearchService, newSingleStaticCalculatorService, newSingleStaticStoreService)

	return service.SearchService{}
}

func NewBondInfoService() service.BondInfoService {
	wire.Build(service.NewBondInfoService, newSingleStaticCalculatorService, newSingleStaticStoreService)

	return service.BondInfoService{}
}
