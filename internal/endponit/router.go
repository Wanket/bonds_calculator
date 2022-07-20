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

	apiGroup := router.fiberApp.Group("/api")

	staticGroup := apiGroup.Group("/static")

	router.search.Configure(staticGroup)
	router.bondInfo.Configure(staticGroup)

	log.Info("Router: configured")
}
