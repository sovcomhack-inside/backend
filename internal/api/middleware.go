package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/utils"
	"github.com/spf13/viper"
)

func (svc *APIService) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie(constants.CookieKeyAuthToken)
		if err != nil {
			return constants.ErrMissingAuthCookie
		}

		token, err := utils.ParseAuthToken(cookie.Value)
		if err != nil {
			return err
		}

		ctx.Set(constants.CtxKeyUserID, token.UserID)

		return next(ctx)
	}
}

func (svc *APIService) OAuthTelegramMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		queryParams := ctx.Request().URL.Query()
		kvs := []string{}
		hash := ""
		for k, v := range queryParams {
			if k == "hash" {
				hash = v[0]
				continue
			}
			kvs = append(kvs, k+"="+v[0])
		}
		sort.Strings(kvs)

		var dataCheckString = ""
		for _, s := range kvs {
			if dataCheckString != "" {
				dataCheckString += "\n"
			}
			dataCheckString += s
		}

		sha256hash := sha256.New()

		telegramToken := viper.GetString("service.telegram_token")
		_, _ = io.WriteString(sha256hash, telegramToken)

		hmachash := hmac.New(sha256.New, sha256hash.Sum(nil))
		_, _ = io.WriteString(hmachash, dataCheckString)

		if hash != hex.EncodeToString(hmachash.Sum(nil)) {
			return constants.ErrHashInvalid
		}

		return next(ctx)
	}
}

func (svc *APIService) AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie(constants.CookieKeySecretToken)
		if err != nil {
			return constants.ErrUnauthorized
		}

		token, err := utils.ParseAuthToken(cookie.Value)
		if err != nil {
			return err
		}

		if token.Secret != viper.GetString(constants.ViperSecretKey) {
			return constants.ErrUnauthorized
		}

		return next(ctx)
	}
}
