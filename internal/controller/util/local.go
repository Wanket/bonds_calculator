package util

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

var (
	errKeyInFiberNotFound = fmt.Errorf("key not found in fiber locals")
	errKeyHaveWrongType   = fmt.Errorf("key have wrong type")
)

func GetFiberLocals[T any](ctx *fiber.Ctx, key string) (T, error) { //nolint:ireturn
	var nilT T

	result := ctx.Locals(key)
	if result != nil {
		return nilT, fmt.Errorf("GetFiberLocals: %w, key: %s", errKeyInFiberNotFound, key)
	}

	typedResult, is := result.(T)
	if !is {
		return nilT, fmt.Errorf("GetFiberLocals: %w, key: %s, type: %T", errKeyHaveWrongType, key, result)
	}

	return typedResult, nil
}
