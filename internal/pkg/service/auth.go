//go:generate mockgen -source=user.go -destination=user_mock.go -package=service
package service

import (
	"context"
	"errors"

	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/logger"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
	"github.com/sovcomhack-inside/internal/pkg/store"
	"github.com/sovcomhack-inside/internal/pkg/utils"
)

type AuthService interface {
	SignupUser(ctx context.Context, request *dto.SignupUserRequest) (*dto.SignupUserResponse, error)
	LoginUser(ctx context.Context, request *dto.LoginUserRequest) (*dto.LoginUserResponse, error)
}

type authServiceImpl struct {
	log   *logger.Logger
	store *store.Registry
}

func (svc *authServiceImpl) SignupUser(ctx context.Context, request *dto.SignupUserRequest) (*dto.SignupUserResponse, error) {
	if _, err := svc.store.UserStore.GetUserByEmail(ctx, request.Email); !errors.Is(err, constants.ErrDBNotFound) {
		if err == nil {
			return nil, constants.ErrEmailAlreadyTaken
		}
		return nil, err
	}

	user := &core.User{
		UserName: request.UserName,
		Email:    request.Email,
	}

	if err := user.UserPassword.Init(request.Password); err != nil {
		return nil, err
	}

	if err := svc.store.UserStore.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	authToken, err := utils.GenerateAuthToken(&utils.AuthTokenWrapper{UserID: user.ID})
	if err != nil {
		return nil, err
	}

	return &dto.SignupUserResponse{AuthToken: authToken}, nil
}

func (svc *authServiceImpl) LoginUser(ctx context.Context, request *dto.LoginUserRequest) (*dto.LoginUserResponse, error) {
	user, err := svc.store.UserStore.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if err := user.UserPassword.Validate(request.Password); err != nil {
		return nil, err
	}

	authToken, err := utils.GenerateAuthToken(&utils.AuthTokenWrapper{UserID: user.ID})
	if err != nil {
		return nil, err
	}

	return &dto.LoginUserResponse{AuthToken: authToken}, nil
}

func NewAuthService(log *logger.Logger, store *store.Registry) AuthService {
	return &authServiceImpl{log: log, store: store}
}
