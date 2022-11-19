package dto

type OAuthTelegramRequest struct {
	ID         int64
	TelegramID int64 `query:"id" validate:"required"`
}
