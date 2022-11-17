package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/utils"
)

func (svc *APIService) AuthMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		cookie := ctx.Cookies(constants.CookieKeyAuthToken)
		if len(cookie) == 0 {
			return constants.ErrMissingAuthCookie
		}

		tw, err := utils.ParseAuthToken(cookie)
		if err != nil {
			return err
		}

		ctx.Request().Header.Set(constants.HeaderKeyUserID, string(tw.UserID))

		return nil
	}
}
