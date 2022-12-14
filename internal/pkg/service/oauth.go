//go:generate mockgen -source=user_test.go -destination=user_mock.go -package=service
package service

import (
	"context"
	"fmt"

	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

func (svc *service) OAuthTelegram(ctx context.Context, request *dto.OAuthTelegramRequest) error {
	if err := svc.store.LinkTelegramID(ctx, request.UserID, request.TelegramID); err != nil {
		return fmt.Errorf("LinkTelegramID: %w", err)
	}
	return nil
}
