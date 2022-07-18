package endponit

import (
	"bonds_calculator/internal/util/marshaling"
	"github.com/gofiber/fiber/v2"
)

func NewFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		JSONEncoder: marshaling.Marshal,
		JSONDecoder: marshaling.Unmarshal,
	})
}
