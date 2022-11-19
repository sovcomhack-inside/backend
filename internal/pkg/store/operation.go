package store

import (
	"context"

	"github.com/sovcomhack-inside/internal/pkg/model/core"
)

type OperationStore interface {
	InsertOperations(ctx context.Context, operations []*core.Operation) error
}

func (s *store) InsertOperations(ctx context.Context, operations []*core.Operation) error {
	query := builder().Insert(tableOperations).
		Columns("purpose", "time", "operation_type", "account_number", "account_amount_cents", "account_amount_currency",
			"original_account_number", "original_amount_cents", "original_amount_currency", "exchange_rate_ratio")
	for _, op := range operations {
		query = query.Values(
			op.Purpose, op.Time, op.OperationType, op.AccountNumber, op.AccountAmountCents, op.AccountAmountCurrency,
			op.OriginalAccountNumber, op.OriginalAmountCents, op.OriginalAmountCurrency, op.ExchangeRateRatio)
	}
	if _, err := s.pool.Execx(ctx, query); err != nil {
		return err
	}

	return nil
}
