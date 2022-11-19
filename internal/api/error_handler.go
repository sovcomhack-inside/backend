package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

func httpErrorHandler(err error, c echo.Context) {
	msg := err.Error()
	for err != nil {
		if ce, ok := err.(*constants.CodedError); ok {
			code := ce.Code()

			_ = c.JSON(code, dto.ErrorResponse{
				Message: msg,
				Code:    code,
			})

			return
		}
		err = errors.Unwrap(err)
	}

	_ = c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
		Message: msg,
		Code:    http.StatusInternalServerError,
	})
}
