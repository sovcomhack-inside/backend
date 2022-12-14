package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
)

type CreateAccountRequest struct {
	UserID   int64  `json:"user_id"`
	Currency string `json:"currency" validate:"required"`
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
	AccountNumber    uuid.UUID       `json:"account_number" validate:"required"`
	DebitAmountCents decimal.Decimal `json:"debit_amount_cents" validate:"required"`
}

type RefillAccountResponse ChangeAccountBalanceResponse

type WithdrawFromAccountRequest struct {
	AccountNumber     uuid.UUID       `json:"account_number" validate:"required"`
	CreditAmountCents decimal.Decimal `json:"credit_amount_cents" validate:"required"`
}

type WithdrawFromAccountResponse ChangeAccountBalanceResponse

type ChangeAccountBalanceResponse struct {
	AccountNumber uuid.UUID       `json:"account_number"`
	OldBalance    decimal.Decimal `json:"old_balance"`
	NewBalance    decimal.Decimal `json:"new_balance"`
	Purpose       string          `json:"purpose"`
}

type MakePurchaseRequest struct {
	AccountNumberFrom  uuid.UUID       `json:"account_number_from" validate:"required"`
	CurrencyFrom       string          `json:"currency_from" validate:"required"`
	DesiredAmountCents decimal.Decimal `json:"desired_amount_cents" validate:"required"`
	AccountNumberTo    uuid.UUID       `json:"account_number_to" validate:"required"`
	CurrencyTo         string          `json:"currency_to" validate:"required"`
}

type MakePurchaseResponse TransferInfo

type MakeSaleRequest struct {
	AccountNumberFrom  uuid.UUID       `json:"account_number_from" validate:"required"`
	CurrencyFrom       string          `json:"currency_from" validate:"required"`
	SellingAmountCents decimal.Decimal `json:"selling_amount_cents" validate:"required"`
	AccountNumberTo    uuid.UUID       `json:"account_number_to" validate:"required"`
	CurrencyTo         string          `json:"currency_to" validate:"required"`
}

type MakeSaleResponse TransferInfo

type TransferInfo struct {
	OldAccountFrom *core.Account `json:"old_account_from"`
	NewAccountFrom *core.Account `json:"new_account_from"`
	OldAccountTo   *core.Account `json:"old_account_to"`
	NewAccountTo   *core.Account `json:"new_account_to"`
	Purpose        string        `json:"purpose"`
}

type TransferRequestDTO struct {
	// AccountFrom ????????, ?? ???????????????? ?????????????????????? ????????????
	AccountFrom uuid.UUID
	// CreditAmountCents ?????????????? ???????????? ?????????? ???? ??????????
	CreditAmountCents decimal.Decimal
	// AccountTo ???????? ???? ?????????????? ?????????????????????? ????????????
	AccountTo uuid.UUID
	// DebitAmountCents ?????????????? ???????????? ?????????????????? ???? ????????
	DebitAmountCents decimal.Decimal
}
