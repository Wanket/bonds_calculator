//go:generate go run github.com/mailru/easyjson/easyjson -no_std_marshalers -lower_camel_case $GOFILE
package controller

import (
	controllerutil "bonds_calculator/internal/controller/util"
	"bonds_calculator/internal/service"
	"bonds_calculator/internal/util"
	"context"
	"errors"
	"github.com/benbjohnson/clock"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type PasswordController struct {
	JWTController
}

//easyjson:json
type changePasswordInfo struct {
	OldPassword string
	NewPassword string
}

func NewPasswordController(
	authService service.IAuthService,
	config util.IGlobalConfig,
	clock clock.Clock,
) *PasswordController {
	return &PasswordController{
		JWTController: NewJWTController(authService, config, clock),
	}
}

func (controller *PasswordController) Configure(router fiber.Router) {
	router.Post("/change-password", controller.ChangePassword)

	log.Info("PasswordController: configured")
}

func (controller *PasswordController) ChangePassword(ctx *fiber.Ctx) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), controller.requestTimeout)
	defer cancel()

	var changePasswordInfo changePasswordInfo
	if err := ctx.BodyParser(&changePasswordInfo); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	username, err := controllerutil.GetFiberLocals[string](ctx, "username")
	if err != nil {
		return controllerutil.SendInternalServerError("PasswordController", ctx, err)
	}

	err = controller.authService.ChangePassword(
		timeoutCtx,
		username,
		changePasswordInfo.OldPassword,
		changePasswordInfo.NewPassword,
	)
	if err != nil {
		if errors.Is(err, service.ErrMismatchedHashAndPassword) {
			return ctx.JSON(controllerutil.BaseResult{
				Status: controllerutil.StatusWrongOldPassword,
			})
		}

		return controllerutil.SendInternalServerError("PasswordController", ctx, err)
	}

	return ctx.JSON(controllerutil.BaseResult{
		Status: controllerutil.StatusOK,
	})
}
