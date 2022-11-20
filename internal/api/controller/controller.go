package controller

import (
	"github.com/sovcomhack-inside/internal/pkg/service"
	"github.com/sovcomhack-inside/internal/pkg/store"
)

type Controller struct {
	store   store.Store
	service service.Service
}

func NewController(store store.Store, service service.Service) *Controller {
	return &Controller{store: store, service: service}
}
