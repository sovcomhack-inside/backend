package dto

type LoginAdminRequest struct {
	Secret string `json:"secret" validate:"required"`
}

type UpdateUserStatusRequest struct {
	ID     int64  `json:"id" validate:"required"`
	Status string `json:"status" validate:"required"`
}

type UpdateUserStatusResponse BasicResponse
