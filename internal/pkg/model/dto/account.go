package dto

import "time"

type CreateAccountRequest struct {
	Currency string `json:"currency"`
}

type CreateAccountResponse struct {
	AccountNumber string    `json:"account_number"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
}
