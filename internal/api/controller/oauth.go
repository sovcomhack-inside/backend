package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
	"github.com/sovcomhack-inside/internal/pkg/utils"
	"github.com/spf13/viper"
)

func (c *Controller) OAuthTelegram(ctx *fiber.Ctx) error {
	request := &dto.OAuthTelegramRequest{}
	if err := Bind(ctx, request, ctx.QueryParser); err != nil {
		return err
	}

	err := c.service.OAuthTelegram(ctx.Context(), request)
	if err != nil {
		return err
	}

	authToken, err := utils.GenerateAuthToken(&utils.AuthTokenWrapper{UserID: request.ID})
	if err != nil {
		return err
	}

	ctx.Cookie(utils.CreateHttpOnlyCookie(constants.CookieKeyAuthToken, authToken, viper.GetInt64(constants.ViperJWTTTLKey)))

	return ctx.Redirect("/", http.StatusPermanentRedirect)
}
