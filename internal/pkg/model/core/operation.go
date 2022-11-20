package core

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OperationType string

const (
	OperationTypeRefill       OperationType = "refill"
	OperationTypeWithdrawal   OperationType = "withdrawal"
	OperationTypeTransfer     OperationType = "transfer"
	OperationTypeSubscription OperationType = "subscription"
)

type Operation struct {
	ID                uuid.UUID        `json:"id" db:"id"`
	Purpose           string           `json:"purpose" db:"purpose"`
	Time              time.Time        `json:"time" db:"time"`
	OperationType     OperationType    `json:"operation_type" db:"operation_type"`
	AccountNumberTo   *uuid.UUID       `json:"account_number_to" db:"account_number_to"`
	AmountCentsTo     *decimal.Decimal `json:"amount_cents_to" db:"amount_cents_to"`
	CurrencyTo        *string          `json:"currency_to" db:"currency_to"`
	AccountNumberFrom *uuid.UUID       `json:"account_number_from" db:"account_number_from"`
	AmountCentsFrom   *decimal.Decimal `json:"amount_cents_from" db:"amount_cents_from"`
	CurrencyFrom      *string          `json:"currency_from" db:"currency_from"`
	ExchangeRateRatio float64          `json:"exchange_rate_ratio" db:"exchange_rate_ratio"`
}
