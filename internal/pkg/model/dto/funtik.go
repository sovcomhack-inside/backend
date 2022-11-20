package dto

import "github.com/google/uuid"

type SubscribeToFuntikRequest struct {
	UserID              int64     `json:"user_id"`
	AccountNumberFrom   uuid.UUID `json:"account_number_from" validate:"required"`
	SubscribePriceCents int64     `json:"subscribe_price_cents" validate:"required"`
}
