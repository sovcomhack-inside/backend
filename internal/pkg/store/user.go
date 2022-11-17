package store

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *core.User) error
	// ListUsers(ctx context.Context) ([]core.User, error)
	// GetUserByID(ctx context.Context, ID string) (*core.User, error)
	GetUserByEmail(ctx context.Context, email string) (*core.User, error)
	// GetUsersByID(ctx context.Context, IDs []string) ([]core.User, error)
	// UpdateUser(ctx context.Context, user *core.User) error
	// DeleteUser(ctx context.Context, user *core.User) error
}

type userStore struct {
	pool Pool
}

func (s *userStore) CreateUser(ctx context.Context, user *core.User) error {
	query := builder().Insert(tableUsers).
		Columns("email", "first_name", "last_name", "password_hash", "password_salt").
		Values(user.Email, user.UserName.First, user.UserName.Last, user.UserPassword.Hash, user.UserPassword.Salt)

	if _, err := s.pool.Execx(ctx, query); err != nil {
		return err
	}

	return nil
}

func (s *userStore) GetUserByEmail(ctx context.Context, email string) (*core.User, error) {
	query := builder().Select("email", "first_name", "last_name", "password_hash", "password_salt").From(tableUsers)
	query = query.Where(squirrel.Eq{"email": email})

	user := &core.User{}
	if err := s.pool.Getx(ctx, user, query); err != nil {
		return nil, wrapErr(err)
	}

	return user, nil
}

func NewUserStore(pool Pool) UserStore {
	return &userStore{pool}
}
