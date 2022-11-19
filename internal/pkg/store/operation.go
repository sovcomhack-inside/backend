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
		Columns("purpose", "time", "operation_type", "receiver_account_number", "receiver_account_amount_cents", "receiver_account_currency",
			"sender_account_number", "sender_account_amount_cents", "sender_account_currency", "currencies_exchange_rate_ratio")
	for _, op := range operations {
		query = query.Values(
			op.Purpose, op.Time, op.OperationType, op.ReceiverAccountNumber, op.ReceiverAmountCents, op.ReceiverAccountCurrency,
			op.SenderAccountNumber, op.SenderAccountAmountCents, op.SenderAccountCurrency, op.ExchangeRateRatio)
	}
	if _, err := s.pool.Execx(ctx, query); err != nil {
		return err
	}

	return nil
}
