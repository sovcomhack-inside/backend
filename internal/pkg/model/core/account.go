package core

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	Number    uuid.UUID       `json:"number" db:"number"`
	UserID    int64           `json:"user_id" db:"user_id"`
	Currency  string          `json:"currency" db:"currency"`
	Balance   decimal.Decimal `json:"balance" db:"balance"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}
