package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
	"github.com/sovcomhack-inside/internal/pkg/service"
	"github.com/sovcomhack-inside/internal/pkg/utils"
	"github.com/spf13/viper"
)

func (c *Controller) SignupUser(ctx *fiber.Ctx) error {
	request := &dto.SignupUserRequest{}
	if err := Bind(ctx, request, ctx.BodyParser); err != nil {
		return err
	}

	response, err := c.service.SignupUser(ctx.Context(), request)
	if err != nil {
		return err
	}

	ctx.Cookie(utils.CreateHttpOnlyCookie(constants.CookieKeyAuthToken, response.AuthToken, viper.GetInt64(constants.ViperJWTTTLKey)))

	return ctx.JSON(response)
}

func (c *Controller) LoginUser(ctx *fiber.Ctx) error {
	request := &dto.LoginUserRequest{}
	if err := Bind(ctx, request, ctx.BodyParser); err != nil {
		return err
	}

	response, err := c.service.LoginUser(ctx.Context(), request)
	if err != nil {
		return err
	}

	ctx.Cookie(utils.CreateHttpOnlyCookie(constants.CookieKeyAuthToken, response.AuthToken, viper.GetInt64(constants.ViperJWTTTLKey)))

	return ctx.JSON(response)
}

func (c *Controller) LogoutUser(ctx *fiber.Ctx) error {
	ctx.ClearCookie(constants.CookieKeyAuthToken)
	return ctx.JSON(&dto.BasicResponse{})
}

func NewAuthController(service service.Service) *Controller {
	return &Controller{service: service}
}
