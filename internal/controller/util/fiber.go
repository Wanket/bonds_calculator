package util

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func SendInternalServerError(prefix string, ctx *fiber.Ctx, err error) error {
	log.WithFields(log.Fields{
		log.ErrorKey: err,
		"url":        ctx.Path(),
	}).Error(prefix + ": InternalServerError")

	return ctx.SendStatus(fiber.StatusInternalServerError)
}
