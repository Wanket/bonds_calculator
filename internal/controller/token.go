package controller

import (
	controllerutil "bonds_calculator/internal/controller/util"
	"bonds_calculator/internal/data/repository"
	"bonds_calculator/internal/service"
	"bonds_calculator/internal/util"
	"context"
	"errors"
	"github.com/benbjohnson/clock"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type TokenController struct {
	JWTController
}

func NewTokenController(
	authService service.IAuthService,
	config util.IGlobalConfig,
	clock clock.Clock,
) *TokenController {
	return &TokenController{
		JWTController: NewJWTController(authService, config, clock),
	}
}

func (controller *TokenController) Configure(router fiber.Router) {
	router.Post("/refresh_token", controller.RefreshToken)

	log.Info("TokenController: configured")
}

func (controller *TokenController) RefreshToken(ctx *fiber.Ctx) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), controller.requestTimeout)
	defer cancel()

	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	jwtToken, err := controllerutil.CheckJwtToken(refreshToken, controller.authService.GetSharedSecret(), controller.clock)
	if err != nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	newToken, err := controller.authService.RefreshToken(timeoutCtx, jwtToken.Bytes)
	if err != nil {
		if errors.Is(err, repository.ErrTokenNotFoundOrExpired) {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		return controllerutil.SendInternalServerError("TokenController", ctx, err)
	}

	return controller.sendJWTCookie(ctx, jwtToken.Username, newToken)
}
