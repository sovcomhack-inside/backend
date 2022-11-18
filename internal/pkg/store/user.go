package store

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
)

var userColumns = []string{"email", "image", "first_name", "last_name", "password_hash", "password_salt"}

func (s *store) CreateUser(ctx context.Context, user *core.User) error {
	query := builder().Insert(tableUsers).
		Columns(userColumns...).
		Values(user.Email, user.Image, user.FirstName, user.LastName, user.UserPassword.Hash, user.UserPassword.Salt)

	if _, err := s.pool.Execx(ctx, query); err != nil {
		return err
	}

	return nil
}

func (s *store) GetUserByEmail(ctx context.Context, email string) (*core.User, error) {
	query := builder().Select(userColumns...).From(tableUsers)
	query = query.Where(squirrel.Eq{"email": email})

	user := &core.User{}
	if err := s.pool.Getx(ctx, user, query); err != nil {
		return nil, wrapErr(err)
	}

	return user, nil
}
