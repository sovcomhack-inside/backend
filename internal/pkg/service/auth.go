//go:generate mockgen -source=user.go -destination=user_mock.go -package=service
package service

import (
	"context"
	"errors"

	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/logger"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
	"github.com/sovcomhack-inside/internal/pkg/utils"
)

func (svc *service) SignupUser(ctx context.Context, request *dto.SignupUserRequest) (*dto.SignupUserResponse, error) {
	if _, err := svc.store.GetUserByEmail(ctx, request.Email); !errors.Is(err, constants.ErrDBNotFound) {
		if err == nil {
			return nil, constants.ErrEmailAlreadyTaken
		}
		return nil, err
	}

	user := &core.User{
		UserName: request.UserName,
		Email:    core.NewNullString(request.Email),
	}
	if err := user.UserPassword.Init(request.Password); err != nil {
		return nil, err
	}

	if err := svc.store.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	mainAccount, err := svc.CreateAccount(ctx, &dto.CreateAccountRequest{
		UserID:   user.ID,
		Currency: "RUB",
	})
	user.MainAccountNumber = mainAccount.Account.Number
	err = svc.store.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	authToken, err := utils.GenerateAuthToken(&utils.AuthTokenWrapper{UserID: user.ID})
	if err != nil {
		return nil, err
	}

	return &dto.SignupUserResponse{User: user, AuthToken: authToken}, nil
}

func (svc *service) LoginUser(ctx context.Context, request *dto.LoginUserRequest) (*dto.LoginUserResponse, error) {
	user, err := svc.store.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if err := user.UserPassword.Validate(request.Password); err != nil {
		return nil, err
	}

	logger.Debugf(ctx, "login: userID: [%v]", user.ID)

	authToken, err := utils.GenerateAuthToken(&utils.AuthTokenWrapper{UserID: user.ID})
	if err != nil {
		return nil, err
	}

	return &dto.LoginUserResponse{User: user, AuthToken: authToken}, nil
}
