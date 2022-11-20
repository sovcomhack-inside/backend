package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
	"github.com/sovcomhack-inside/internal/pkg/store"
)

func (c *Controller) ListOperations(ctx echo.Context) error {
	var request = &dto.ListOperationsRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	opts := &store.SearchOperationsOpts{
		AccountNumbersIn: request.AccountNumbersIn,
		OperationTypesIn: request.OperationTypesIn,
	}
	operations, err := c.service.ListOperations(ctx.Request().Context(), opts)
	if err != nil {
		return err
	}

	resp := &dto.ListOperationsResponse{
		Operations: operations,
	}
	return ctx.JSON(http.StatusOK, resp)
}
