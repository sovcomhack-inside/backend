package dto

import (
	"github.com/sovcomhack-inside/internal/pkg/model/core"
)

type CreateAccountRequest struct {
	Currency string `json:"currency"`
}

type CreateAccountResponse struct {
	Account core.Account `json:"account"`
}

type ListUserAccountsRequest struct {
	UserID int64 `json:"user_id"`
}

type ListUserAccountResponse struct {
	Accounts []core.Account `json:"accounts"`
}
