package endponit

import (
	"bonds_calculator/internal/util/marshaling"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"
	log "github.com/sirupsen/logrus"
)

func NewFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: marshaling.Marshal,
		JSONDecoder: marshaling.Unmarshal,

		ErrorHandler: errorHandler,
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: panicHandler,
	}))

	return app
}

func panicHandler(ctx *fiber.Ctx, err any) {
	log.WithFields(log.Fields{
		log.ErrorKey: err,
		"url":        ctx.Path(),
	}).Error("Fiber: got panic on request")
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return ctx.Status(code).SendString(utils.StatusMessage(code))
}
