package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

func (c *Controller) CreateAccount(ctx echo.Context) error {
	request := &dto.CreateAccountRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if ctx.Get(constants.CtxKeyUserID) == nil {
		return fmt.Errorf("can't resolve user_id")
	}
	request.UserID = ctx.Get(constants.CtxKeyUserID).(int64)
	response, err := c.service.CreateAccount(ctx.Request().Context(), request)
	if err != nil {
		return fmt.Errorf("account controller internal error: %w", err)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *Controller) ListUserAccounts(ctx echo.Context) error {
	request := &dto.ListUserAccountsRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}
	if ctx.Get(constants.CtxKeyUserID) == nil {
		return fmt.Errorf("can't resolve user_id")
	}
	request.UserID = ctx.Get(constants.CtxKeyUserID).(int64)
	response, err := c.service.ListUserAccounts(ctx.Request().Context(), request)
	if err != nil {
		return fmt.Errorf("account controller internal error: %w", err)

	}
	return ctx.JSON(http.StatusOK, response)
}

func (c *Controller) RefillAccount(ctx echo.Context) error {
	request := &dto.RefillAccountRequest{}

	if err := ctx.Bind(request); err != nil {
		return err
	}
	if request.DebitAmountCents <= 0 {
		return constants.ErrNegativeDebit
	}
	response, err := c.service.RefillAccount(ctx.Request().Context(), request)
	if err != nil {
		return fmt.Errorf("account controller internal error: %w", err)
	}
	return ctx.JSON(http.StatusOK, response)
}

func (c *Controller) WithdrawFromAccount(ctx echo.Context) error {
	request := &dto.WithdrawFromAccountRequest{}

	if err := ctx.Bind(request); err != nil {
		return err
	}
	if request.CreditAmountCents <= 0 {
		return constants.ErrNegativeCredit
	}
	response, err := c.service.WithdrawFromAccount(ctx.Request().Context(), request)
	if err != nil {
		if errors.Is(err, constants.ErrNotEnoughMoney) {
			return constants.ErrNotEnoughMoney
		}
		return fmt.Errorf("account controller internal error: %w", err)
	}
	return ctx.JSON(http.StatusOK, response)
}
