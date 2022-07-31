package endponit

import (
	"bonds_calculator/internal/controller"
	"bonds_calculator/internal/midleware"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type Router struct {
	fiberApp *fiber.App

	search   *controller.SearchController
	bondInfo *controller.BondInfoController
	auth     *controller.AuthController
	token    *controller.TokenController
	password *controller.PasswordController

	authMiddleware midleware.AuthMiddleware
}

func NewRouter(
	app *fiber.App,
	search *controller.SearchController,
	bondInfo *controller.BondInfoController,
	auth *controller.AuthController,
	token *controller.TokenController,
	authMiddleware midleware.AuthMiddleware,
	password *controller.PasswordController,
) *Router {
	return &Router{
		fiberApp: app,

		search:   search,
		bondInfo: bondInfo,
		auth:     auth,
		token:    token,
		password: password,

		authMiddleware: authMiddleware,
	}
}

func (router *Router) Configure() {
	log.Info("Router: configuring")

	router.configureStatic()

	apiGroup := router.fiberApp.Group("/api")

	apiStaticGroup := apiGroup.Group("/static")

	router.search.Configure(apiStaticGroup)
	router.bondInfo.Configure(apiStaticGroup)

	apiUserGroup := apiGroup.Group("/user")
	apiUserGroup.Use(router.authMiddleware)

	router.password.Configure(apiUserGroup)

	apiSessionUserGroup := apiGroup.Group("/session")

	router.auth.Configure(apiSessionUserGroup)
	router.token.Configure(apiSessionUserGroup)

	log.Info("Router: configured")
}

func (router *Router) configureStatic() {
	router.fiberApp.Static("/", "./public")

	router.fiberApp.Static("/page/*", "./public/index.html")
}
