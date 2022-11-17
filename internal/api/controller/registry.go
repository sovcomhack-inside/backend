package controller

import (
	"github.com/sovcomhack-inside/internal/pkg/logger"
	"github.com/sovcomhack-inside/internal/pkg/service"
	"github.com/sovcomhack-inside/internal/pkg/store"
)

type Registry struct {
	AuthController *AuthController
}

func NewRegistry(log *logger.Logger, dbRegistry *store.Registry) (*Registry, error) {
	serviceRegistry := service.NewRegistry(log, dbRegistry)

	return &Registry{
		AuthController: NewAuthController(log, serviceRegistry),
	}, nil
}
