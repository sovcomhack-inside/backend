package store

import (
	"context"

	"github.com/sovcomhack-inside/internal/pkg/model/core"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
	"github.com/sovcomhack-inside/internal/pkg/store/xpgx"
)

type Tx = xpgx.Tx
type Pool = xpgx.Pool

type Store interface {
	UserStore
	AccountStore
	OperationStore
}

type store struct {
	pool Pool
}

func NewStore(pool Pool) Store {
	return &store{pool}
}

type UserStore interface {
	CreateUser(ctx context.Context, user *core.User) error
	// ListUsers(ctx context.Context) ([]core.User, error)
	// GetUserByID(ctx context.Context, ID string) (*core.User, error)
	GetUserByEmail(ctx context.Context, email string) (*core.User, error)
	// GetUsersByID(ctx context.Context, IDs []string) ([]core.User, error)
	// UpdateUser(ctx context.Context, user *core.User) error
	// DeleteUser(ctx context.Context, user *core.User) error
	UpdateUserStatus(ctx context.Context, id int64, status core.UserStatus) error
	GetUserStatus(ctx context.Context, id int64) (string, error)
	ListUsers(ctx context.Context, request *dto.ListUsersRequest) ([]core.User, error)
	LinkTelegramID(ctx context.Context, id, telegramID int64) error
	GetTelegramID(ctx context.Context, id int64) (int64, error)
}
