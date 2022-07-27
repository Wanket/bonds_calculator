package controller

import (
	"bonds_calculator/internal/controller"
	"bonds_calculator/internal/endponit"
	"github.com/gofiber/fiber/v2"
)

func CreateAppAndRegistryController(url string, controller controller.IController) *fiber.App {
	app := endponit.NewFiberApp()
	controller.Configure(app.Group(url))

	return app
}
