package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

func (c *Controller) GetUser(ctx echo.Context) error {
	id, err := GetUserIDFromCtx(ctx)
	if err != nil {
		return err
	}

	user, err := c.store.GetUserByID(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &dto.GetUserResponse{User: user})
}
