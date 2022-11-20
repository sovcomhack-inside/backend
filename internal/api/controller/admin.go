package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
	"github.com/sovcomhack-inside/internal/pkg/utils"
	"github.com/spf13/viper"
)

func (c *Controller) LoginAdmin(ctx echo.Context) error {
	request := &dto.LoginAdminRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if request.Secret != viper.GetString(constants.ViperSecretKey) {
		return constants.ErrUnauthorized
	}

	authToken, err := utils.GenerateAuthToken(&utils.AuthTokenWrapper{Secret: request.Secret})
	if err != nil {
		return err
	}

	ctx.SetCookie(utils.CreateHttpOnlyCookie(constants.CookieKeySecretToken, authToken, viper.GetInt64(constants.ViperJWTTTLKey)))

	return ctx.JSON(http.StatusOK, &dto.BasicResponse{})
}

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

func (c *Controller) ListUsers(ctx echo.Context) error {
	request := &dto.ListUsersRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	users, err := c.store.ListUsers(ctx.Request().Context(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &dto.ListUsersResponse{Users: users})
}
