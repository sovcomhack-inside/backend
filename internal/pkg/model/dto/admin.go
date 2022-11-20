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
	Status  string   `json:"status"`
	EmailIn []string `json:"email_in"`
}

type ListUsersResponse struct {
	Users []core.User `json:"users"`
}
