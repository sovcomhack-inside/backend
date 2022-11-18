package dto

import "github.com/sovcomhack-inside/internal/pkg/model/core"

type SignupUserRequest struct {
	core.UserName
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignupUserResponse struct {
	*core.User
	AuthToken string `json:"auth_token"`
}

type LoginUserRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponse struct {
	*core.User
	AuthToken string `json:"auth_token"`
}
