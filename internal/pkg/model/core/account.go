package core

import "github.com/google/uuid"

type Account struct {
	Number   uuid.UUID `db:"number"`
	UserID   int64     `db:"user_id"`
	Currency string    `db:"currency"`
	Cents    int64     `db:"cents"`
}
