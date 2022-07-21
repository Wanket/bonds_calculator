package controller

import "github.com/gofiber/fiber/v2"

type IController interface {
	Configure(router fiber.Router)
}
