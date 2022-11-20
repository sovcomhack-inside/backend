package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
	"github.com/sovcomhack-inside/internal/pkg/store"
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
	MakeTransfer(ctx context.Context, reqDTO *dto.TransferRequestDTO, exchangeRateRatio decimal.Decimal) (*dto.MakePurchaseResponse, error)
}

func (svc *service) CreateAccount(ctx context.Context, req *dto.CreateAccountRequest) (*dto.CreateAccountResponse, error) {
	account := &core.Account{
		Number:   uuid.New(),
		UserID:   req.UserID,
		Currency: req.Currency,
		Balance:  decimal.NewFromInt(0),
	}

	err := svc.store.CreateAccount(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("account store error: %w", err)
	}
	return &dto.CreateAccountResponse{
		Account: *account,
	}, nil
}

func (svc *service) ListUserAccounts(ctx context.Context, req *dto.ListUserAccountsRequest) (*dto.ListUserAccountResponse, error) {
	accounts, err := svc.store.SearchUserAccounts(ctx, store.SearchAccountsOpts{UserID: req.UserID})
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
			Purpose:           refillPurpose,
			OperationType:     core.OperationTypeRefill,
			AccountNumberTo:   &req.AccountNumber,
			AmountCentsTo:     &req.DebitAmountCents,
			CurrencyTo:        &account.Currency,
			AccountNumberFrom: nil,
			AmountCentsFrom:   nil,
			CurrencyFrom:      nil,
			ExchangeRateRatio: decimal.NewFromInt(1),
		},
	}
	err = svc.store.InsertOperations(ctx, operations)
	if err != nil {
		return nil, fmt.Errorf("operations store err: %w", err)
	}
	return &dto.RefillAccountResponse{
		AccountNumber: account.Number,
		OldBalance:    account.Balance.Sub(req.DebitAmountCents),
		NewBalance:    account.Balance,
		Purpose:       refillPurpose,
	}, nil
}

func (svc *service) WithdrawFromAccount(ctx context.Context, req *dto.WithdrawFromAccountRequest) (*dto.WithdrawFromAccountResponse, error) {
	account, err := svc.store.UpdateAccountBalance(ctx, req.AccountNumber, req.CreditAmountCents.Neg())
	if err != nil {
		return nil, fmt.Errorf("accounts store error: %w", err)
	}

	operations := []*core.Operation{
		{
			Purpose:           withdrawalPurpose,
			OperationType:     core.OperationTypeWithdrawal,
			AccountNumberTo:   nil,
			AmountCentsTo:     nil,
			CurrencyTo:        nil,
			AccountNumberFrom: &req.AccountNumber,
			AmountCentsFrom:   lo.ToPtr(req.CreditAmountCents.Neg()),
			CurrencyFrom:      &account.Currency,
			ExchangeRateRatio: decimal.NewFromInt(1),
		},
	}
	err = svc.store.InsertOperations(ctx, operations)
	if err != nil {
		return nil, fmt.Errorf("operations store err: %w", err)
	}
	return &dto.WithdrawFromAccountResponse{
		AccountNumber: account.Number,
		OldBalance:    account.Balance.Add(req.CreditAmountCents),
		NewBalance:    account.Balance,
		Purpose:       withdrawalPurpose,
	}, nil
}

// MakeTransfer перевести со счета на счет
func (svc *service) MakeTransfer(ctx context.Context, reqDTO *dto.TransferRequestDTO, exchangeRateRatio decimal.Decimal) (*dto.MakePurchaseResponse, error) {
	accFrom, accTo, err := svc.store.TransferMoney(ctx, reqDTO)
	if err != nil {
		return nil, fmt.Errorf("accounts store error: %w", err)
	}
	accFromBefore := &core.Account{
		Number:    accFrom.Number,
		UserID:    accFrom.UserID,
		Currency:  accFrom.Currency,
		Balance:   accFrom.Balance.Add(reqDTO.CreditAmountCents),
		CreatedAt: accFrom.CreatedAt,
	}
	accToBefore := &core.Account{
		Number:    accTo.Number,
		UserID:    accTo.UserID,
		Currency:  accTo.Currency,
		Balance:   accTo.Balance.Sub(reqDTO.DebitAmountCents),
		CreatedAt: accTo.CreatedAt,
	}

	var purpose = fmt.Sprintf("Перевод со счета %s (%s) на счет %s (%s)", accFrom.Number.String(), accFrom.Currency, accTo.Number.String(), accTo.Currency)
	operations := []*core.Operation{
		{
			Purpose:           purpose,
			OperationType:     core.OperationTypeTransfer,
			AccountNumberTo:   &accTo.Number,
			AmountCentsTo:     &reqDTO.DebitAmountCents,
			CurrencyTo:        &accTo.Currency,
			AccountNumberFrom: &accFrom.Number,
			AmountCentsFrom:   &reqDTO.CreditAmountCents,
			CurrencyFrom:      &accFrom.Currency,
			ExchangeRateRatio: exchangeRateRatio,
		},
	}
	err = svc.store.InsertOperations(ctx, operations)
	if err != nil {
		return nil, fmt.Errorf("operations store err: %w", err)
	}
	return &dto.MakePurchaseResponse{
		OldAccountFrom: accFromBefore,
		NewAccountFrom: accFrom,
		OldAccountTo:   accToBefore,
		NewAccountTo:   accTo,
		Purpose:        purpose,
	}, nil
}
