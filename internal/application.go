package internal

import (
	"bonds_calculator/internal/endponit"
	"bonds_calculator/internal/util"
	"fmt"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Application struct {
	fiberApp *fiber.App

	config util.IGlobalConfig
}

func NewApplication(app *fiber.App, router *endponit.Router, config util.IGlobalConfig) *Application {
	router.Configure()

	return &Application{
		fiberApp: app,

		config: config,
	}
}

func (application *Application) Run() error {
	log.Infof("Application: starting listening on port %d", application.config.HTTPPort())

	return fmt.Errorf("fiber: %w", application.fiberApp.Listen(":"+strconv.Itoa(application.config.HTTPPort())))
}
