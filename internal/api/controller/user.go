package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

func (c *Controller) UpdateUserStatus(ctx echo.Context) error {
	request := &dto.UpdateUserStatusRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	response, err := c.service.UpdateUserStatus(ctx.Request().Context(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}
