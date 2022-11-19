package dto

type OAuthTelegramRequest struct {
	ID        int64  `query:"id" validate:"required"`
	FirstName string `query:"first_name" validate:"required"`
	LastName  string `query:"last_name" validate:"required"`
	PhotoURL  string `query:"photo_url"`
}
