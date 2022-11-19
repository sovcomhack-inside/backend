package store

import (
	"context"

	"github.com/sovcomhack-inside/internal/pkg/model/core"
	"github.com/sovcomhack-inside/internal/pkg/store/xpgx"
)

type Tx = xpgx.Tx
type Pool = xpgx.Pool

type Store interface {
	UserStore
	AccountStore
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
}
