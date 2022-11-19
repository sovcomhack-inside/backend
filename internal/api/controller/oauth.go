package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
	"github.com/sovcomhack-inside/internal/pkg/utils"
	"github.com/spf13/viper"
)

func (c *Controller) OAuthTelegram(ctx echo.Context) error {
	request := &dto.OAuthTelegramRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	err := c.service.OAuthTelegram(ctx.Request().Context(), request)
	if err != nil {
		return err
	}

	authToken, err := utils.GenerateAuthToken(&utils.AuthTokenWrapper{UserID: request.ID})
	if err != nil {
		return err
	}

	ctx.SetCookie(utils.CreateHttpOnlyCookie(constants.CookieKeyAuthToken, authToken, viper.GetInt64(constants.ViperJWTTTLKey)))

	return ctx.Redirect(http.StatusPermanentRedirect, "/")
}
