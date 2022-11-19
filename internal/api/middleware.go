package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/utils"
	"github.com/spf13/viper"
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

func (svc *APIService) OAuthTelegramMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		queryParams := ctx.Request().URI().QueryArgs()
		kvs := []string{}
		hash := ""
		queryParams.VisitAll(func(k []byte, v []byte) {
			if string(k) == "hash" {
				hash = string(v)
			} else {
				kvs = append(kvs, string(k)+"="+string(v))
			}
		})
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

		return nil
	}
}
