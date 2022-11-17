package controller

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/logger"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
	"github.com/sovcomhack-inside/internal/pkg/service"
	"github.com/sovcomhack-inside/internal/pkg/utils"
	"github.com/spf13/viper"
)

type AuthController struct {
	log      *logger.Logger
	registry *service.Registry
}

func (c *AuthController) SignupUser(ctx *fiber.Ctx) error {
	request := &dto.SignupUserRequest{}
	if err := Bind(ctx, request, ctx.BodyParser); err != nil {
		return err
	}

	response, err := c.registry.AuthService.SignupUser(context.Background(), request)
	if err != nil {
		return err
	}

	ctx.Cookie(utils.CreateCookie(constants.CookieKeyAuthToken, response.AuthToken, viper.GetInt64(constants.ViperJWTTTLKey)))

	return ctx.JSON(nil)
}

func (c *AuthController) LoginUser(ctx *fiber.Ctx) error {
	request := &dto.LoginUserRequest{}
	if err := Bind(ctx, request, ctx.BodyParser); err != nil {
		return err
	}

	response, err := c.registry.AuthService.LoginUser(context.Background(), request)
	if err != nil {
		return err
	}

	ctx.Cookie(utils.CreateCookie(constants.CookieKeyAuthToken, response.AuthToken, viper.GetInt64(constants.ViperJWTTTLKey)))

	return ctx.JSON(nil)
}

func (c *AuthController) LogoutUser(ctx *fiber.Ctx) error {
	ctx.ClearCookie(constants.CookieKeyAuthToken)
	return ctx.JSON(&dto.BasicResponse{})
}

func NewAuthController(log *logger.Logger, registry *service.Registry) *AuthController {
	return &AuthController{log: log, registry: registry}
}
