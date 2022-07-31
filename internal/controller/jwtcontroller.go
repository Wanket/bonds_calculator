package controller

import (
	controllerutil "bonds_calculator/internal/controller/util"
	"bonds_calculator/internal/service"
	"bonds_calculator/internal/util"
	"github.com/benbjohnson/clock"
	"github.com/gofiber/fiber/v2"
	"time"
)

type JWTController struct {
	authService service.IAuthService

	requestTimeout time.Duration
	accessTokenTTL time.Duration

	clock clock.Clock
}

func NewJWTController(authService service.IAuthService, config util.IGlobalConfig, clock clock.Clock) JWTController {
	return JWTController{
		authService: authService,

		requestTimeout: config.GetDefaultRequestTimeout(),
		accessTokenTTL: config.GetAccessTokenTTL(),

		clock: clock,
	}
}

func (controller *JWTController) sendJWTCookie(ctx *fiber.Ctx, username string, refreshToken service.Token) error {
	sharedSecret := controller.authService.GetSharedSecret()

	refreshCookie, err := controllerutil.CreateJWTCookie(username, "refresh_token", sharedSecret, refreshToken)
	if err != nil {
		return controllerutil.SendInternalServerError("JWTController", ctx, err)
	}

	accessSecret, err := util.GenerateSecretString(util.TokenSize)
	if err != nil {
		return controllerutil.SendInternalServerError("JWTController", ctx, err)
	}

	accessToken := service.Token{
		Bytes:          accessSecret,
		ExpirationTime: controller.clock.Now().Add(controller.accessTokenTTL),
	}

	accessCookie, err := controllerutil.CreateJWTCookie(username, "access_token", sharedSecret, accessToken)
	if err != nil {
		return controllerutil.SendInternalServerError("JWTController", ctx, err)
	}

	ctx.Cookie(&refreshCookie)
	ctx.Cookie(&accessCookie)

	return ctx.JSON(controllerutil.BaseResult{
		Status: controllerutil.ResultOK,
	})
}
