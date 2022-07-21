package controller

import (
	"bonds_calculator/internal/controller"
	"github.com/gofiber/fiber/v2"
)

func createAppAndRegistryController(url string, controller controller.IController) *fiber.App {
	app := fiber.New()
	controller.Configure(app.Group(url))

	return app
}
