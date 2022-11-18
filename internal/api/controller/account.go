package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sovcomhack-inside/internal/pkg/logger"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
	"github.com/sovcomhack-inside/internal/pkg/service"
)

type AccountController struct {
	log            *logger.Logger
	accountService service.AccountService
}

func NewAccountController(log *logger.Logger, accountService service.AccountService) *AccountController {
	return &AccountController{
		log:            log,
		accountService: accountService,
	}
}

func (c *AccountController) CreateAccount(ctx *fiber.Ctx) error {
	request := &dto.CreateAccountRequest{}
	if err := Bind(ctx, request, ctx.BodyParser); err != nil {
		return err
	}

	response, err := c.accountService.CreateAccount(ctx.Context(), request)
	if err != nil {
		return err
	}
	return ctx.JSON(response)
}
