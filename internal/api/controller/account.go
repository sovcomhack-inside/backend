package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
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
	if request.DebitAmountCents.LessThanOrEqual(decimal.NewFromInt(0)) {
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
	if request.CreditAmountCents.LessThanOrEqual(decimal.NewFromInt(0)) {
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

func (c *Controller) MakePurchase(ctx echo.Context) error {
	request := &dto.MakePurchaseRequest{}

	if err := ctx.Bind(request); err != nil {
		return err
	}
	if request.DesiredAmountCents.LessThanOrEqual(decimal.NewFromInt(0)) {
		return constants.ErrNegativeDebit
	}
	rate := getRate(request.CurrencyTo, request.CurrencyFrom)
	reqDTO := &dto.TransferRequestDTO{
		AccountFrom:       request.AccountNumberFrom,
		CreditAmountCents: decimal.NewFromFloat(rate).Mul(request.DesiredAmountCents),
		AccountTo:         request.AccountNumberTo,
		DebitAmountCents:  request.DesiredAmountCents,
	}
	response, err := c.service.MakeTransfer(ctx.Request().Context(), reqDTO, rate)
	if err != nil {
		if errors.Is(err, constants.ErrNotEnoughMoney) {
			return constants.ErrNotEnoughMoney
		}
		return fmt.Errorf("account controller internal error: %w", err)
	}
	return ctx.JSON(http.StatusOK, response)
}

func getRate(fromCurrency, toCurrency string) float64 {
	var rates = map[string]map[string]float64{
		"RUB": {
			"EUR": 0.0123,
			"USD": 0.01,
		},
		"USD": {
			"RUB": 65.21523452,
			"EUR": 1.0010200012026,
		},
		"EUR": {
			"RUB": 62.952736,
			"USD": 0.98182,
		},
	}
	return rates[fromCurrency][toCurrency]
}
