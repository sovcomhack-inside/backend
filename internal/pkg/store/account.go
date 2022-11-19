package store

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
)

type AccountStore interface {
	CreateAccount(ctx context.Context, account *core.Account) error
	SearchUserAccounts(ctx context.Context, userID int64) ([]core.Account, error)
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
	query := builder().Select("number", "user_id", "currency", "cents", "created_at").
		From(tableAccounts).
		Where(squirrel.Eq{"user_id": userID})

	var accounts []core.Account

	if err := s.pool.Selectx(ctx, &accounts, query); err != nil {
		return nil, err
	}

	return accounts, nil
}
