package controller

import (
	"github.com/sovcomhack-inside/internal/pkg/logger"
	"github.com/sovcomhack-inside/internal/pkg/service"
)

type Controller struct {
	log     *logger.Logger
	service service.Service
}

func NewController(log *logger.Logger, service service.Service) *Controller {
	return &Controller{log: log, service: service}
}
