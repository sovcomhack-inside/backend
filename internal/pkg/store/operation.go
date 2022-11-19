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
		Columns("purpose", "time", "operation_type", "account_number_to", "amount_cents_to", "currency_to",
			"account_number_from", "amount_cents_from", "currency_from", "currencies_exchange_rate_ratio")
	for _, op := range operations {
		query = query.Values(
			op.Purpose, op.Time, op.OperationType, op.AccountNumberTo, op.AmountCentsTo, op.CurrencyTo,
			op.AccountNumberFrom, op.AmountCentsFrom, op.CurrencyFrom, op.ExchangeRateRatio)
	}
	if _, err := s.pool.Execx(ctx, query); err != nil {
		return err
	}

	return nil
}
