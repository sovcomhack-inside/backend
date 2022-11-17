package dto

import "github.com/sovcomhack-inside/internal/pkg/model/core"

type SignupUserRequest struct {
	core.UserName
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignupUserResponse AuthTokenResponse

type LoginUserRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponse AuthTokenResponse
