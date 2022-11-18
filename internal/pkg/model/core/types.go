package core

import "database/sql"

type NullString = sql.NullString

func NewNullString(s string) NullString {
	return sql.NullString{String: s, Valid: len(s) > 0}
}
