package endponit

import (
	"bonds_calculator/internal/controller"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type Router struct {
	fiberApp *fiber.App

	search *controller.SearchController
}

func NewRouter(app *fiber.App, search *controller.SearchController) *Router {
	return &Router{
		fiberApp: app,
		search:   search,
	}
}

func (router *Router) Configure() {
	log.Info("Router: configuring")

	apiGroup := router.fiberApp.Group("/api")

	router.search.Configure(apiGroup)

	log.Info("Router: configured")
}
