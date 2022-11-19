package dto

import "github.com/sovcomhack-inside/internal/pkg/model/core"

type LoginAdminRequest struct {
	Secret string `json:"secret" validate:"required"`
}

type UpdateUserStatusRequest struct {
	ID     int64  `json:"id" validate:"required"`
	Status string `json:"status" validate:"required"`
}

type UpdateUserStatusResponse BasicResponse

type ListUsersRequest struct {
	Status string `query:"status" validate:"required"`
}

type ListUsersResponse struct {
	Users []core.User `json:"users"`
}
