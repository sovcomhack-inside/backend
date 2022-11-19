package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
	"github.com/sovcomhack-inside/internal/pkg/store/xpgx"
)

type AccountStore interface {
	CreateAccount(ctx context.Context, account *core.Account) error
	SearchUserAccounts(ctx context.Context, userID int64) ([]core.Account, error)
	UpdateAccountBalance(ctx context.Context, accountNumber uuid.UUID, debitAmount int64) (acc *core.Account, txErr error)
}

func (s *store) CreateAccount(ctx context.Context, account *core.Account) error {
	query := builder().Insert(tableAccounts).
		Columns("number", "user_id", "currency").
		Values(account.Number, account.UserID, account.Currency).
		Suffix("ON CONFLICT (user_id, currency) DO NOTHING RETURNING created_at")

	if err := s.pool.QueryxRow(ctx, query).Scan(&account.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return constants.AccountAlreadyExists
		}
		return err
	}

	return nil
}

func (s *store) SearchUserAccounts(ctx context.Context, userID int64) ([]core.Account, error) {
	query := builder().Select("number", "user_id", "currency", "balance", "created_at").
		From(tableAccounts).
		Where(squirrel.Eq{"user_id": userID})

	var accounts []core.Account

	if err := s.pool.Selectx(ctx, &accounts, query); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *store) UpdateAccountBalance(ctx context.Context, accountNumber uuid.UUID, deltaAmountCents int64) (acc *core.Account, txErr error) {
	txErr = s.withTx(ctx, func(ctx context.Context, tx Tx) error {
		var err error

		acc, err = getAccount(ctx, accountNumber, tx)
		if err != nil {
			return fmt.Errorf("getAccount err: %w", err)
		}
		if acc.Balance+deltaAmountCents < 0 {
			return constants.ErrNotEnoughMoney
		}

		newBalance, err := updateAccountBalance(ctx, accountNumber, acc.Balance, deltaAmountCents, tx)
		if err != nil {
			return fmt.Errorf("updateAccountBalance err: %w", err)
		}
		acc.Balance = newBalance
		return nil
	})
	return
}

func getAccount(ctx context.Context, accountNumber uuid.UUID, executor xpgx.Executor) (*core.Account, error) {
	var account core.Account

	query := builder().Select("number", "user_id", "currency", "balance", "created_at").
		From(tableAccounts).
		Where(squirrel.Eq{"number": accountNumber})

	err := executor.Getx(ctx, &account, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("account with this number not found: %s", accountNumber.String())
		}
		return nil, err
	}
	return &account, err
}

func updateAccountBalance(ctx context.Context, accountNumber uuid.UUID, balance, delta int64, executor xpgx.Executor) (int64, error) {
	query := builder().Update(tableAccounts).
		Set("balance", balance+delta).
		Where(squirrel.Eq{"number": accountNumber})

	if _, err := executor.Execx(ctx, query); err != nil {
		return 0, err
	}
	return balance + delta, nil
}
