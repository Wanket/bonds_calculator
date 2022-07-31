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

func newAuthMiddleware() midleware.AuthMiddleware {
	wire.Build(midleware.NewAuth, getSingleAuthService, clock.New)

	return nil
}

func newApplication() *Application {
	wire.Build(
		NewApplication,
		newAuthMiddleware,
		endponit.NewRouter,
		endponit.NewFiberApp,
		getSingleSearchController,
		getSingleBondInfoController,
		getSingleConfig,
		getSingleAuthController,
		getSingleTokenController,
		getSinglePasswordController,
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

func newAuthController() *controller.AuthController {
	wire.Build(controller.NewAuthController, clock.New, getSingleAuthService, getSingleConfig)

	return &controller.AuthController{}
}

func newTokenController() *controller.TokenController {
	wire.Build(controller.NewTokenController, clock.New, getSingleAuthService, getSingleConfig)

	return &controller.TokenController{}
}

func newAuthService() *service.AuthService {
	wire.Build(
		service.NewAuthService,
		clock.New,
		getSingleTokenRepository,
		getSingleUserRepositoryBuilder,
		getSingleConfig,
	)

	return &service.AuthService{}
}

func newTokenRepository() *repository.TokenRepository {
	wire.Build(repository.NewTokenRepository, getSingleRedisConnection)

	return &repository.TokenRepository{}
}

func newUserRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, getSingleDBConnection)

	return &repository.UserRepository{}
}

func newRedisConnection() *connection.RedisConnection {
	wire.Build(connection.NewRedisConnection, getSingleConfig)

	return &connection.RedisConnection{}
}

func newDBConnection() *connection.DBConnection {
	wire.Build(connection.NewDBConnection, getSingleConfig)

	return &connection.DBConnection{}
}

func newPasswordController() *controller.PasswordController {
	wire.Build(controller.NewPasswordController, getSingleAuthService, getSingleConfig, clock.New)

	return &controller.PasswordController{}
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
	authService             = lazy.New(newAuthService)

	searchController   = lazy.New(newSearchController)
	bondInfoController = lazy.New(newBondInfoController)
	authController     = lazy.New(newAuthController)
	tokenController    = lazy.New(newTokenController)
	passwordController = lazy.New(newPasswordController)

	tokenRepository = lazy.New(newTokenRepository)
	userRepository  = lazy.New(newUserRepository)

	redisConnection = lazy.New(newRedisConnection)
	dbConnection    = lazy.New(newDBConnection)

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

func getSingleAuthController() *controller.AuthController {
	return authController.Value()
}

func getSingleTokenController() *controller.TokenController {
	return tokenController.Value()
}

func getSingleAuthService() service.IAuthService {
	return authService.Value()
}

func getSingleTokenRepository() repository.ITokenRepository {
	return tokenRepository.Value()
}

func getSingleUserRepository() repository.IUserRepository {
	return userRepository.Value()
}

func getSingleUserRepositoryBuilder() datautil.IRepositoryBuilder[repository.IUserRepository] {
	return userRepository.Value()
}

func getSingleRedisConnection() *connection.RedisConnection {
	return redisConnection.Value()
}

func getSingleDBConnection() *connection.DBConnection {
	return dbConnection.Value()
}

func getSinglePasswordController() *controller.PasswordController {
	return passwordController.Value()
}
