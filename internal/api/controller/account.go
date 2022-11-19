package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

func (c *Controller) CreateAccount(ctx echo.Context) error {
	request := &dto.CreateAccountRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	response, err := c.service.CreateAccount(ctx.Request().Context(), request)
	if err != nil {
		return fmt.Errorf("account controller internal error: %w", err)
	}

	return ctx.JSON(http.StatusOK, response)
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

func (c *Controller) WithdrawFromAccount(ctx *fiber.Ctx) error {
	request := &dto.WithdrawFromAccountRequest{}

	if err := Bind(ctx, request, ctx.BodyParser); err != nil {
		return err
	}
	if request.CreditAmountCents <= 0 {
		return constants.ErrNegativeCredit
	}
	response, err := c.service.WithdrawFromAccount(ctx.Context(), request)
	if err != nil {
		if errors.Is(err, constants.ErrNotEnoughMoney) {
			return constants.ErrNotEnoughMoney
		}
		return fmt.Errorf("account controller internal error: %w", err)
	}
	return ctx.JSON(response)
}
