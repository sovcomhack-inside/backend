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

	currencies := c.service.ListCurrencies(ctx.Request().Context(), request.Code)

	return ctx.JSON(http.StatusOK, &dto.ListCurrenciesResponse{Currencies: currencies})
}
