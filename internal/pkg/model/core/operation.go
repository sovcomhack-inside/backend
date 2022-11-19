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
	ID                       uuid.UUID     `json:"id" db:"id"`
	Purpose                  string        `json:"purpose" db:"purpose"`
	Time                     time.Time     `json:"time" db:"time"`
	OperationType            OperationType `json:"operation_type" db:"operation_type"`
	ReceiverAccountNumber    *uuid.UUID    `json:"receiver_account_number" db:"receiver_account_number"`
	ReceiverAmountCents      *int64        `json:"receiver_account_amount_cents" db:"receiver_account_amount_cents"`
	ReceiverAccountCurrency  *string       `json:"receiver_account_currency" db:"receiver_account_currency"`
	SenderAccountNumber      *uuid.UUID    `json:"sender_account_number" db:"sender_account_number"`
	SenderAccountAmountCents *int64        `json:"sender_account_amount_cents" db:"sender_account_amount_cents"`
	SenderAccountCurrency    *string       `json:"sender_account_currency" db:"sender_account_currency"`
	ExchangeRateRatio        float64       `json:"currencies_exchange_rate_ratio" db:"currencies_exchange_rate_ratio"`
}
