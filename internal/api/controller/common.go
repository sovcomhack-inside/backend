package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sovcomhack-inside/internal/pkg/constants"
)

type Parser func(out interface{}) error

var validate = validator.New()

func ValidateStruct(i interface{}) error {
	if err := validate.Struct(i); err != nil {
		var violations []string
		for _, ve := range err.(validator.ValidationErrors) {
			violations = append(violations, fmt.Sprintf("[%s] didn't satisfy [%s]", ve.StructNamespace(), ve.Tag()))
		}
		return constants.NewCodedError(strings.Join(violations, ","), http.StatusBadRequest)
	}
	return nil
}

func Bind(ctx *fiber.Ctx, out interface{}, parsers ...Parser) error {
	for _, parser := range parsers {
		parser(out)
	}

	if err := ValidateStruct(out); err != nil {
		return err
	}

	return nil
}
