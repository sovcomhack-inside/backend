package utils

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func SetValue(ctx *fiber.Ctx, key any, value any) {
	ctx.SetUserContext(context.WithValue(ctx.UserContext(), key, value))
}

func GetValue(ctx *fiber.Ctx, key any) any {
	return ctx.UserContext().Value(key)
}
