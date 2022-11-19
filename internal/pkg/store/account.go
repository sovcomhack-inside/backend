package store

import (
	"context"

	"github.com/sovcomhack-inside/internal/pkg/model/core"
)

type AccountStore interface {
	CreateAccount(ctx context.Context, account *core.Account) error
}

func (s *store) CreateAccount(ctx context.Context, account *core.Account) error {
	query := builder().Insert(tableAccounts).
		Columns("number", "user_id", "currency").
		Values(account.Number, account.UserID, account.Currency)

	if _, err := s.pool.Execx(ctx, query); err != nil {
		return err
	}

	return nil
}
