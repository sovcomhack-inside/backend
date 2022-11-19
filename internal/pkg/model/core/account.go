package core

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Number    uuid.UUID `json:"number" db:"number"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Currency  string    `json:"currency" db:"currency"`
	Balance   int64     `json:"balance" db:"balance"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
