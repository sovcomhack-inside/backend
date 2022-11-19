package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

const (
	refillPurpose     = "Пополнение счета с банковской карты"
	withdrawalPurpose = "Вывод со счета на банковскую карту"
)

type AccountService interface {
	CreateAccount(context.Context, *dto.CreateAccountRequest) (*dto.CreateAccountResponse, error)
	ListUserAccounts(context.Context, *dto.ListUserAccountsRequest) (*dto.ListUserAccountResponse, error)
	RefillAccount(context.Context, *dto.RefillAccountRequest) (*dto.RefillAccountResponse, error)
	WithdrawFromAccount(context.Context, *dto.WithdrawFromAccountRequest) (*dto.WithdrawFromAccountResponse, error)
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
		Balance:  0,
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
		return nil, fmt.Errorf("accounts store error: %w", err)
	}
	return &dto.ListUserAccountResponse{Accounts: accounts}, nil
}

func (svc *service) RefillAccount(ctx context.Context, req *dto.RefillAccountRequest) (*dto.RefillAccountResponse, error) {
	account, err := svc.store.UpdateAccountBalance(ctx, req.AccountNumber, req.DebitAmountCents)
	if err != nil {
		return nil, fmt.Errorf("accounts store error: %w", err)
	}

	operations := []*core.Operation{
		{
			Purpose:                  refillPurpose,
			OperationType:            core.OperationTypeRefill,
			ReceiverAccountNumber:    &req.AccountNumber,
			ReceiverAmountCents:      &req.DebitAmountCents,
			ReceiverAccountCurrency:  &account.Currency,
			SenderAccountNumber:      nil,
			SenderAccountAmountCents: nil,
			SenderAccountCurrency:    nil,
			ExchangeRateRatio:        1,
		},
	}
	err = svc.store.InsertOperations(ctx, operations)
	if err != nil {
		return nil, fmt.Errorf("operations store err: %w", err)
	}
	return &dto.RefillAccountResponse{
		OldBalance: account.Balance - req.DebitAmountCents,
		NewBalance: account.Balance,
		Purpose:    refillPurpose,
	}, nil
}

func (svc *service) WithdrawFromAccount(ctx context.Context, req *dto.WithdrawFromAccountRequest) (*dto.WithdrawFromAccountResponse, error) {
	account, err := svc.store.UpdateAccountBalance(ctx, req.AccountNumber, -req.CreditAmountCents)
	if err != nil {
		return nil, fmt.Errorf("accounts store error: %w", err)
	}

	operations := []*core.Operation{
		{
			Purpose:                  withdrawalPurpose,
			OperationType:            core.OperationTypeWithdrawal,
			ReceiverAccountNumber:    nil,
			ReceiverAmountCents:      nil,
			ReceiverAccountCurrency:  nil,
			SenderAccountNumber:      &req.AccountNumber,
			SenderAccountAmountCents: lo.ToPtr(-req.CreditAmountCents),
			SenderAccountCurrency:    &account.Currency,
			ExchangeRateRatio:        1,
		},
	}
	err = svc.store.InsertOperations(ctx, operations)
	if err != nil {
		return nil, fmt.Errorf("operations store err: %w", err)
	}
	return &dto.WithdrawFromAccountResponse{
		OldBalance: account.Balance + req.CreditAmountCents,
		NewBalance: account.Balance,
		Purpose:    withdrawalPurpose,
	}, nil
}
