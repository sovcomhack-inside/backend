package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

func (c *Controller) CreateAccount(ctx *fiber.Ctx) error {
	request := &dto.CreateAccountRequest{}
	if err := Bind(ctx, request, ctx.BodyParser); err != nil {
		return err
	}

	response, err := c.service.CreateAccount(ctx.Context(), request)
	if err != nil {
		return err
	}
	return ctx.JSON(response)
}

func (c *Controller) ListUserAccounts(ctx *fiber.Ctx) error {
	request := &dto.ListUserAccountsRequest{}
	if err := Bind(ctx, request, ctx.BodyParser); err != nil {
		return err
	}

	response, err := c.service.ListUserAccounts(ctx.Context(), request)
	if err != nil {
		return err
	}
	return ctx.JSON(response)
}
