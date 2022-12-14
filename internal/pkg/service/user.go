package service

import (
	"context"

	"github.com/sovcomhack-inside/internal/pkg/model/core"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

func (svc *service) UpdateUserStatus(ctx context.Context, request *dto.UpdateUserStatusRequest) (*dto.UpdateUserStatusResponse, error) {
	if err := svc.store.UpdateUserStatus(ctx, request.ID, core.UserStatus(request.Status)); err != nil {
		return nil, err
	}

	return &dto.UpdateUserStatusResponse{}, svc.store.UpdateUserStatus(ctx, request.ID, core.UserStatus(request.Status))
}
