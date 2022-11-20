package service

import (
	"context"

	"github.com/sovcomhack-inside/internal/pkg/model/core"
)

type OperationService interface {
	ListOperations(context.Context, []*core.Operation) error
}

func (svc *service) ListOperations(ctx context.Context, operations []*core.Operation) error {
	return nil
}
