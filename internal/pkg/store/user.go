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

func (s *store) ListUsersInStatus(ctx context.Context, status core.UserStatus) ([]core.User, error) {
	users := []core.User{}
	const query = `
	SELECT id, email, image, first_name, last_name FROM users WHERE id IN (
		SELECT id FROM users_statuses WHERE status = $1
	)
	`
	if err := s.pool.Select(ctx, &users, query, status); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *store) LinkTelegramID(ctx context.Context, id, telegramID int64) error {
	const query = `
	INSERT INTO users_telegram_id COLUMNS (id, user_id) VALUES ($1, $2)
	`
	if _, err := s.pool.Exec(ctx, query, id, telegramID); err != nil {
		return err
	}
	return nil
}

func (s *store) GetTelegramID(ctx context.Context, id int64) (telegramID int64, err error) {
	const query = "SELECT telegram_id FROM users_telegram_id WHERE id = $1"
	err = s.pool.Get(ctx, &telegramID, query, id)
	return
}
