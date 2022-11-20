package dto

import (
	"github.com/google/uuid"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
)

type ListOperationsRequest struct {
	AccountNumbersIn []uuid.UUID          `json:"account_numbers_in"`
	OperationTypesIn []core.OperationType `json:"operation_types_in"`
}

type ListOperationsResponse struct {
	Operations []*core.Operation `json:"operations"`
}
