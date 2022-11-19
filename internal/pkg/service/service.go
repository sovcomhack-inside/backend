package service

import (
	"context"

	"github.com/sovcomhack-inside/internal/pkg/model/dto"
	"github.com/sovcomhack-inside/internal/pkg/store"
)

type Service interface {
	UserService
	AuthService
	OAuthService
	AccountService
}

type service struct {
	store store.Store
}

func NewService(store store.Store) Service {
	return &service{store}
}

type UserService interface {
	UpdateUserStatus(ctx context.Context, request *dto.UpdateUserStatusRequest) (*dto.UpdateUserStatusResponse, error)
}

type AuthService interface {
	SignupUser(ctx context.Context, request *dto.SignupUserRequest) (*dto.SignupUserResponse, error)
	LoginUser(ctx context.Context, request *dto.LoginUserRequest) (*dto.LoginUserResponse, error)
}

type OAuthService interface {
	OAuthTelegram(ctx context.Context, request *dto.OAuthTelegramRequest) error
}
