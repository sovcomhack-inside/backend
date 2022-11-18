package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sovcomhack-inside/internal/pkg/logger"
	"github.com/sovcomhack-inside/internal/pkg/store/xpgx"
)

func (s *store) withTx(ctx context.Context, fn func(ctx context.Context, tx Tx) error) (err error) {
	pgxTx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}

	tx := xpgx.NewTx(pgxTx)

	defer func() {
		if p := recover(); p != nil {
			rollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			rollbackTx(ctx, tx)
		} else {
			if txErr := tx.Commit(ctx); txErr != nil {
				err = fmt.Errorf("cannot commit transaction: %w", err)
			}
		}
	}()

	return fn(ctx, tx)
}

func rollbackTx(ctx context.Context, tx Tx) {
	if err := tx.Rollback(ctx); err != nil {
		if !errors.Is(err, pgx.ErrTxClosed) {
			logger.Errorf(ctx, "tx.Rollback failed: %v", err)
		}
	}
}
