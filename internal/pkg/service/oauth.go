//go:generate mockgen -source=user_test.go -destination=user_mock.go -package=service
package service

import (
	"context"
	"fmt"

	"github.com/sovcomhack-inside/internal/pkg/model/core"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

func (svc *service) OAuthTelegram(ctx context.Context, request *dto.OAuthTelegramRequest) error {
	user := &core.User{
		ID:       request.ID,
		UserName: core.UserName{FirstName: request.FirstName, LastName: request.LastName},
		Image:    request.PhotoURL,
	}

	if err := svc.store.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("CreateUser: %w", err)
	}

	return nil
}
