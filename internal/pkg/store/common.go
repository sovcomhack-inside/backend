package store

import (
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/sovcomhack-inside/internal/pkg/constants"
)

const (
	tableUsers         = "users"
	tableAccounts      = "accounts"
	tableOperations    = "operations"
	tableUsersStatuses = "users_statuses"
)

var mapping = map[error]error{pgx.ErrNoRows: constants.ErrDBNotFound}

func wrapErr(err error) error {
	for k, v := range mapping {
		if errors.Is(err, k) {
			return v
		}
	}
	return err
}

// builder возвращает squirrel SQL Builder обьект.
func builder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
