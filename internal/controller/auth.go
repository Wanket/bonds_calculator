//go:generate go run github.com/mailru/easyjson/easyjson -no_std_marshalers -lower_camel_case $GOFILE
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

type AuthController struct {
	JWTController
}

//easyjson:json
type loginInfo struct {
	Login    string
	Password string
}

func NewAuthController(authService service.IAuthService, config util.IGlobalConfig, clock clock.Clock) *AuthController {
	return &AuthController{
		JWTController: NewJWTController(authService, config, clock),
	}
}

func (controller *AuthController) Configure(router fiber.Router) {
	router.Post("/login", controller.Login)
	router.Post("/register", controller.Register)

	log.Info("AuthController: configured")
}

func (controller *AuthController) Register(ctx *fiber.Ctx) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), controller.requestTimeout)
	defer cancel()

	var loginInfo loginInfo
	if err := ctx.BodyParser(&loginInfo); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	token, err := controller.authService.CreateUser(timeoutCtx, loginInfo.Login, loginInfo.Password)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return ctx.JSON(controllerutil.BaseResult{
				Status: controllerutil.StatusUserAlreadyExists,
			})
		}

		return controllerutil.SendInternalServerError("AuthController", ctx, err)
	}

	return controller.sendJWTCookie(ctx, loginInfo.Login, token)
}

func (controller *AuthController) Login(ctx *fiber.Ctx) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), controller.requestTimeout)
	defer cancel()

	var loginInfo loginInfo

	if err := ctx.BodyParser(&loginInfo); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	token, err := controller.authService.LoginUser(timeoutCtx, loginInfo.Login, loginInfo.Password)
	if err != nil {
		if errors.Is(err, service.ErrMismatchedHashAndPassword) {
			return ctx.JSON(controllerutil.BaseResult{
				Status: controllerutil.StatusWrongLoginOrPassword,
			})
		}

		return controllerutil.SendInternalServerError("AuthController", ctx, err)
	}

	return controller.sendJWTCookie(ctx, loginInfo.Login, token)
}
