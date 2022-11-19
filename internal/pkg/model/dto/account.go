package dto

import (
	"github.com/google/uuid"
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

type RefillAccountRequest struct {
	AccountNumber    uuid.UUID `json:"account_number"`
	DebitAmountCents int64     `json:"debit_amount_cents"`
}

type RefillAccountResponse ChangeAccountBalanceResponse

type WithdrawFromAccountRequest struct {
	AccountNumber     uuid.UUID `json:"account_number"`
	CreditAmountCents int64     `json:"credit_amount_cents"`
}

type WithdrawFromAccountResponse ChangeAccountBalanceResponse

type ChangeAccountBalanceResponse struct {
	OldBalance int64  `json:"old_balance"`
	NewBalance int64  `json:"new_balance"`
	Purpose    string `json:"purpose"`
}
