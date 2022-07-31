//go:generate go run github.com/mailru/easyjson/easyjson -lower_camel_case $GOFILE
package util

import (
	"bonds_calculator/internal/service"
	"bonds_calculator/internal/util"
	"errors"
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/gofiber/fiber/v2"
	"github.com/kataras/jwt"
)

var errTokenExpired = errors.New("token expired")

//easyjson:json
type JWTToken struct {
	service.Token

	Username string
}

func CreateJWTCookie(username string, cookieName string, jwtKey []byte, token service.Token) (fiber.Cookie, error) {
	signedToken, err := jwt.Sign(jwt.HS256, jwtKey, JWTToken{
		Token: token,

		Username: username,
	})
	if err != nil {
		return fiber.Cookie{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return fiber.Cookie{
		Name:     cookieName,
		Value:    string(signedToken),
		Expires:  token.ExpirationTime,
		Secure:   true,
		HTTPOnly: true,
	}, nil
}

func CheckJwtToken(stringToken string, jwtKey []byte, clock clock.Clock) (JWTToken, error) {
	jwtToken, err := jwt.Verify(jwt.HS256, jwtKey, util.StringToBytes(stringToken))
	if err != nil {
		return JWTToken{}, fmt.Errorf("failed to verify token: %w", err)
	}

	var token JWTToken
	if err := jwtToken.Claims(&token); err != nil {
		return JWTToken{}, fmt.Errorf("failed to parse token: %w", err)
	}

	if token.ExpirationTime.Before(clock.Now()) {
		return JWTToken{}, errTokenExpired
	}

	return token, nil
}
