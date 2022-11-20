package dto

type OAuthTelegramRequest struct {
	UserID     int64 `param:"user_id" validate:"required"`
	TelegramID int64 `query:"id" validate:"required"`
}
