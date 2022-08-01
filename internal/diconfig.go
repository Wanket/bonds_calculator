//go:build wireinject

//go:generate go run github.com/google/wire/cmd/wire

package internal

import (
	"bonds_calculator/internal/api"
	"bonds_calculator/internal/controller"
	"bonds_calculator/internal/data/connection"
	"bonds_calculator/internal/data/repository"
	datautil "bonds_calculator/internal/data/util"
	"bonds_calculator/internal/endponit"
	"bonds_calculator/internal/midleware"
	"bonds_calculator/internal/service"
	"bonds_calculator/internal/util"
	"github.com/benbjohnson/clock"
	"github.com/google/wire"
)

func CreateApp() *Application {
	wire.Build(
		NewApplication,

		midleware.NewAuthMiddleware,

		endponit.NewRouter,
		endponit.NewFiberApp,

		controller.NewSearchController,
		controller.NewBondInfoController,
		controller.NewAuthController,
		controller.NewTokenController,
		controller.NewPasswordController,

		service.NewSearchService,
		wire.Bind(new(service.ISearchService), new(*service.SearchService)),
		service.NewBondInfoService,
		wire.Bind(new(service.IBondInfoService), new(*service.BondInfoService)),
		service.NewAuthService,
		wire.Bind(new(service.IAuthService), new(*service.AuthService)),
		service.NewStaticStoreService,
		wire.Bind(new(service.IStaticStoreService), new(*service.StaticStoreService)),
		service.NewTimerService,
		wire.Bind(new(service.ITimerService), new(*service.TimerService)),
		service.NewStaticCalculatorService,
		wire.Bind(new(service.IStaticCalculatorService), new(*service.StaticCalculatorService)),

		repository.NewUserRepository,
		//wire.Bind(new(repository.IUserRepository), new(*repository.UserRepository)),
		wire.Bind(new(datautil.IRepositoryBuilder[repository.IUserRepository]), new(*repository.UserRepository)),
		repository.NewTokenRepository,
		wire.Bind(new(repository.ITokenRepository), new(*repository.TokenRepository)),

		connection.NewRedisConnection,
		connection.NewDBConnection,

		clock.New,

		defaultMoexClientProvider,
		wire.Bind(new(api.IMoexClient), new(*api.MoexClient)),
		util.NewTimeHelper,
		wire.Bind(new(util.ITimeHelper), new(*util.TimeHelper)),
		util.NewGlobalConfig,
		wire.Bind(new(util.IGlobalConfig), new(*util.GlobalConfig)),
	)

	return &Application{}
}

func defaultMoexClientProvider(config util.IGlobalConfig) *api.MoexClient {
	return api.NewMoexClient(config.MoexClientQueueSize())
}
