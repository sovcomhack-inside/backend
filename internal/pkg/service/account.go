package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

type AccountService interface {
	CreateAccount(context.Context, *dto.CreateAccountRequest) (*dto.CreateAccountResponse, error)
}

func (s *service) CreateAccount(ctx context.Context, accountDTO *dto.CreateAccountRequest) (*dto.CreateAccountResponse, error) {
	userId, err := strconv.ParseInt(ctx.Value(constants.CtxKeyUserID{}).(string), 10, 64)
	if err != nil {
		return nil, err
	}
	account := &core.Account{
		Number:   uuid.New(),
		UserID:   userId,
		Currency: accountDTO.Currency,
		Cents:    0,
	}

	err = s.store.CreateAccount(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("account store error: %w", err)
	}
	return &dto.CreateAccountResponse{
		AccountNumber: account.Number.String(),
		Currency:      account.Currency,
	}, nil
}
