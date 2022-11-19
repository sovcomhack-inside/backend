package store

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
)

var userColumns = []string{"id", "email", "image", "first_name", "last_name", "password_hash", "password_salt"}

func (s *store) CreateUser(ctx context.Context, user *core.User) error {
	return s.withTx(ctx, func(ctx context.Context, tx Tx) error {
		query := builder().Insert(tableUsers).
			Columns(userColumns[1:]...).
			Values(user.Email, user.Image, user.FirstName, user.LastName, user.UserPassword.Hash, user.UserPassword.Salt).
			Suffix("RETURNING id")

		if err := tx.Getx(ctx, user, query); err != nil {
			return err
		}

		query = builder().Insert(tableUsersStatuses).Columns("id").Values(user.ID)
		if _, err := tx.Execx(ctx, query); err != nil {
			return err
		}

		return nil
	})
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

func (s *store) UpdateUserStatus(ctx context.Context, id int64, status core.UserStatus) error {
	query := builder().Update(tableUsersStatuses).Set("status", status).Where(squirrel.Eq{"id": id})
	if _, err := s.pool.Execx(ctx, query); err != nil {
		return err
	}
	return nil
}

func (s *store) GetUserStatus(ctx context.Context, id int64) (string, error) {
	var status string
	query := builder().Select("status").From(tableUsersStatuses).Where(squirrel.Eq{"id": id})
	if err := s.pool.Getx(ctx, &status, query); err != nil {
		return "", err
	}
	return status, nil
}
