package endponit

import (
	"bonds_calculator/internal/controller"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type Router struct {
	fiberApp *fiber.App

	search   *controller.SearchController
	bondInfo *controller.BondInfoController
}

func NewRouter(app *fiber.App, search *controller.SearchController, bondInfo *controller.BondInfoController) *Router {
	return &Router{
		fiberApp: app,

		search:   search,
		bondInfo: bondInfo,
	}
}

func (router *Router) Configure() {
	log.Info("Router: configuring")

	router.configureStatic()

	apiGroup := router.fiberApp.Group("/api")

	apiStaticGroup := apiGroup.Group("/static")

	router.search.Configure(apiStaticGroup)
	router.bondInfo.Configure(apiStaticGroup)

	log.Info("Router: configured")
}

func (router *Router) configureStatic() {
	router.fiberApp.Static("/", "./public")

	router.fiberApp.Static("/page/*", "./public/index.html")
}
