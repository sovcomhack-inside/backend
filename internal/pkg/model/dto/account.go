package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
)

type CreateAccountRequest struct {
	UserID   int64  `json:"user_id"`
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
	AccountNumber    uuid.UUID       `json:"account_number"`
	DebitAmountCents decimal.Decimal `json:"debit_amount_cents"`
}

type RefillAccountResponse ChangeAccountBalanceResponse

type WithdrawFromAccountRequest struct {
	AccountNumber     uuid.UUID       `json:"account_number"`
	CreditAmountCents decimal.Decimal `json:"credit_amount_cents"`
}

type WithdrawFromAccountResponse ChangeAccountBalanceResponse

type ChangeAccountBalanceResponse struct {
	AccountNumber uuid.UUID       `json:"account_number"`
	OldBalance    decimal.Decimal `json:"old_balance"`
	NewBalance    decimal.Decimal `json:"new_balance"`
	Purpose       string          `json:"purpose"`
}

type MakePurchaseRequest struct {
	AccountNumberFrom  uuid.UUID       `json:"account_number_from"`
	CurrencyFrom       string          `json:"currency_from"`
	DesiredAmountCents decimal.Decimal `json:"desired_amount_cents"`
	AccountNumberTo    uuid.UUID       `json:"account_number_to"`
	CurrencyTo         string          `json:"currency_to"`
}

type MakePurchaseResponse TransferInfo

type MakeSaleRequest struct {
	AccountNumberFrom  uuid.UUID       `json:"account_number_from"`
	CurrencyFrom       string          `json:"currency_from"`
	SellingAmountCents decimal.Decimal `json:"selling_amount_cents"`
	AccountNumberTo    uuid.UUID       `json:"account_number_to"`
	CurrencyTo         string          `json:"currency_to"`
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
	// AccountFrom счет, с которого переводятся деньги
	AccountFrom uuid.UUID
	// CreditAmountCents сколько копеек снять со счета
	CreditAmountCents decimal.Decimal
	// AccountTo счет на который переводятся деньги
	AccountTo uuid.UUID
	// DebitAmountCents сколько копеек перевести на счет
	DebitAmountCents decimal.Decimal
}
