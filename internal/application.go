package internal

import (
	"bonds_calculator/internal/endponit"
	"bonds_calculator/internal/util"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Application struct {
	fiberApp *fiber.App
}

func NewApplication(app *fiber.App, router *endponit.Router) *Application {
	router.Configure()

	return &Application{
		fiberApp: app,
	}
}

func (application *Application) Run() error {
	log.Infof("Application: starting listening on port %d", util.GetGlobalConfig().HttpPort())

	return application.fiberApp.Listen(":" + strconv.Itoa(util.GetGlobalConfig().HttpPort()))
}
