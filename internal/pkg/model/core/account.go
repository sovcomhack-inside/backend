package core

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Number    uuid.UUID `db:"number"`
	UserID    int64     `db:"user_id"`
	Currency  string    `db:"currency"`
	Cents     int64     `db:"cents"`
	CreatedAt time.Time `db:"created_at"`
}
