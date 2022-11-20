package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

func (c *Controller) ListCurrencies(ctx echo.Context) error {
	var request = &dto.ListCurrenciesRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	currencies, err := c.service.ListCurrencies(ctx.Request().Context(), request.Code)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &dto.ListCurrenciesResponse{Currencies: currencies})
}

func (c *Controller) GetCurrencyData(ctx echo.Context) error {
	var request = &dto.GetCurrencyDataRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	resp, err := c.service.GetCurrencyData(ctx.Request().Context(), request.Code, request.BaseCurrencyCode, request.DaysNumber)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, resp)
}
