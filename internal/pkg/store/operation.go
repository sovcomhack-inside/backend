package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
)

type SearchOperationsOpts struct {
	AccountNumbersIn []uuid.UUID
	OperationTypesIn []core.OperationType
}

type OperationStore interface {
	InsertOperations(ctx context.Context, operations []*core.Operation) error
	SearchOperations(ctx context.Context, opts *SearchOperationsOpts) ([]*core.Operation, error)
}

func (s *store) InsertOperations(ctx context.Context, operations []*core.Operation) error {
	query := builder().Insert(tableOperations).
		Columns("purpose", "operation_type", "account_number_to", "amount_cents_to", "currency_to",
			"account_number_from", "amount_cents_from", "currency_from", "exchange_rate_ratio")
	for _, op := range operations {
		query = query.Values(
			op.Purpose, op.OperationType, op.AccountNumberTo, op.AmountCentsTo, op.CurrencyTo,
			op.AccountNumberFrom, op.AmountCentsFrom, op.CurrencyFrom, op.ExchangeRateRatio)
	}
	if _, err := s.pool.Execx(ctx, query); err != nil {
		return err
	}

	return nil
}

func (s *store) SearchOperations(ctx context.Context, opts *SearchOperationsOpts) ([]*core.Operation, error) {
	query := builder().Select("id", "purpose", "time", "operation_type", "account_number_to", "amount_cents_to", "currency_to",
		"account_number_from", "amount_cents_from", "currency_from", "exchange_rate_ratio").
		From(tableOperations)
	if len(opts.AccountNumbersIn) > 0 {
		numbersStr := strings.Join(
			lo.Map(
				opts.AccountNumbersIn,
				func(number uuid.UUID, _ int) string {
					return fmt.Sprintf("'%s'", number.String())
				}),
			", ",
		)
		query = query.Where(fmt.Sprintf("account_number_to in(%s) OR account_number_from in (%s)", numbersStr, numbersStr))
	}

	if len(opts.OperationTypesIn) > 0 {
		query = query.Where(squirrel.Eq{"operation_type": opts.OperationTypesIn})
	}
	query = query.OrderBy("time DESC")
	var operations []*core.Operation
	if err := s.pool.Selectx(ctx, &operations, query); err != nil {
		return nil, err
	}

	return operations, nil
}
