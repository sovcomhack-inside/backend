package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

func (c *Controller) CreateAccount(ctx *fiber.Ctx) error {
	request := &dto.CreateAccountRequest{}
	if err := Bind(ctx, request, ctx.BodyParser); err != nil {
		return err
	}

	response, err := c.service.CreateAccount(ctx.Context(), request)
	if err != nil {
		return fmt.Errorf("account controller internal error: %w", err)

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
		return fmt.Errorf("account controller internal error: %w", err)

	}
	return ctx.JSON(response)
}

func (c *Controller) RefillAccount(ctx *fiber.Ctx) error {
	request := &dto.RefillAccountRequest{}

	if err := Bind(ctx, request, ctx.BodyParser); err != nil {
		return err
	}
	if request.DebitAmountCents <= 0 {
		return constants.ErrNegativeDebit
	}
	response, err := c.service.RefillAccount(ctx.Context(), request)
	if err != nil {
		return fmt.Errorf("account controller internal error: %w", err)
	}
	return ctx.JSON(response)
}
