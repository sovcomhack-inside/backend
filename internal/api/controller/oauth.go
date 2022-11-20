package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

func (c *Controller) OAuthTelegram(ctx echo.Context) error {
	request := &dto.OAuthTelegramRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if err := c.service.OAuthTelegram(ctx.Request().Context(), request); err != nil {
		return err
	}

	return ctx.Redirect(http.StatusPermanentRedirect, "/")
}
