package controller

import (
	"github.com/sovcomhack-inside/internal/pkg/service"
)

type Controller struct {
	service service.Service
}

func NewController(service service.Service) *Controller {
	return &Controller{service: service}
}
