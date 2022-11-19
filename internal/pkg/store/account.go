package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
)

type AccountStore interface {
	CreateAccount(ctx context.Context, account *core.Account) error
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
