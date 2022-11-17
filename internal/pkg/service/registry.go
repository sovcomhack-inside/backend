package service

import (
	"github.com/sovcomhack-inside/internal/pkg/logger"
	"github.com/sovcomhack-inside/internal/pkg/store"
)

type Registry struct {
	AuthService AuthService
}

func NewRegistry(log *logger.Logger, store *store.Registry) *Registry {
	return &Registry{
		AuthService: NewAuthService(log, store),
	}
}
