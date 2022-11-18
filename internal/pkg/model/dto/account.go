package dto

type CreateAccountRequest struct {
	Currency string `json:"currency"`
}

type CreateAccountResponse struct {
	AccountNumber string `json:"account_number"`
	Currency      string `json:"currency"`
}
