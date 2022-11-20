package dto

type OAuthTelegramRequest struct {
	UserID     int64 `query:"param" validate:"required"`
	TelegramID int64 `query:"id" validate:"required"`
}
