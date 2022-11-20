package service

import (
	"context"

	"github.com/sovcomhack-inside/internal/pkg/model/core"
	"github.com/sovcomhack-inside/internal/pkg/store"
)

type OperationService interface {
	ListOperations(ctx context.Context, opts *store.SearchOperationsOpts) ([]*core.Operation, error)
}

func (svc *service) ListOperations(ctx context.Context, opts *store.SearchOperationsOpts) ([]*core.Operation, error) {
	operations, err := svc.store.SearchOperations(ctx, opts)
	if err != nil {
		return nil, err
	}
	return operations, nil
}
