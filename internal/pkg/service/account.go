package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

type AccountService interface {
	CreateAccount(context.Context, *dto.CreateAccountRequest) (*dto.CreateAccountResponse, error)
	ListUserAccounts(ctx context.Context, req *dto.ListUserAccountsRequest) (*dto.ListUserAccountResponse, error)
}

func (svc *service) CreateAccount(ctx context.Context, accountDTO *dto.CreateAccountRequest) (*dto.CreateAccountResponse, error) {
	//userId, err := strconv.ParseInt(ctx.Value(constants.CtxKeyUserID{}).(string), 10, 64)
	//if err != nil {
	//	return nil, err
	//}
	var userId int64
	var err error

	userId = 1
	account := &core.Account{
		Number:   uuid.New(),
		UserID:   userId,
		Currency: accountDTO.Currency,
		Cents:    0,
	}

	err = svc.store.CreateAccount(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("account store error: %w", err)
	}
	return &dto.CreateAccountResponse{
		Account: *account,
	}, nil
}

func (svc *service) ListUserAccounts(ctx context.Context, req *dto.ListUserAccountsRequest) (*dto.ListUserAccountResponse, error) {
	//userId, err := strconv.ParseInt(ctx.Value(constants.CtxKeyUserID{}).(string), 10, 64)
	//if err != nil {
	//	return nil, err
	//}
	var userId int64
	var err error

	userId = 1
	accounts, err := svc.store.SearchUserAccounts(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("account store error: %w", err)
	}
	return &dto.ListUserAccountResponse{Accounts: accounts}, nil
}
