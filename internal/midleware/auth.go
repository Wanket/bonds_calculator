package midleware

import (
	controllerutil "bonds_calculator/internal/controller/util"
	"bonds_calculator/internal/service"
	"github.com/benbjohnson/clock"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware = fiber.Handler

func NewAuth(auth service.IAuthService, clock clock.Clock) AuthMiddleware {
	jwtKey := auth.GetSharedSecret()

	return func(ctx *fiber.Ctx) error {
		accessToken := ctx.Cookies("access_token")
		if accessToken == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		jwtToken, err := controllerutil.CheckJwtToken(accessToken, jwtKey, clock)
		if err != nil {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		ctx.Locals("username", jwtToken.Username)

		return nil
	}
}
