package core

import (
	"time"

	"github.com/google/uuid"
)

type OperationType string

const (
	OperationTypeRefill           OperationType = "refill"
	OperationTypeWithdrawal       OperationType = "withdrawal"
	OperationTypeTransferIncoming OperationType = "transfer_incoming"
	OperationTypeTransferOutgoing OperationType = "transfer_outgoing"
)

type Operation struct {
	ID                     uuid.UUID     `json:"id" db:"id"`
	Purpose                string        `json:"purpose" db:"purpose"`
	Time                   time.Time     `json:"time" db:"time"`
	OperationType          OperationType `json:"operation_type" db:"operation_type"`
	AccountNumber          *uuid.UUID    `json:"account_number" db:"account_number"`
	AccountAmountCents     *int64        `json:"account_amount_cents" db:"account_amount_cents"`
	AccountAmountCurrency  *string       `json:"account_amount_currency" db:"account_amount_currency"`
	OriginalAccountNumber  *uuid.UUID    `json:"original_account_number" db:"original_account_number"`
	OriginalAmountCents    *int64        `json:"original_amount_cents" db:"original_amount_cents"`
	OriginalAmountCurrency *string       `json:"original_amount_currency" db:"original_amount_currency"`
	ExchangeRateRatio      float64       `json:"exchange_rate_ratio" db:"exchange_rate_ratio"`
}
