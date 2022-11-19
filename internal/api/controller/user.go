package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

func (c *Controller) UpdateUserStatus(ctx *fiber.Ctx) error {
	request := &dto.UpdateUserStatusRequest{}
	if err := Bind(ctx, request, ctx.BodyParser); err != nil {
		return err
	}

	response, err := c.service.UpdateUserStatus(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}
