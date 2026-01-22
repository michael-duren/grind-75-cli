package utils

import (
	"database/sql"
	"time"
)

func CoerceFromNullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func CoerceFromNullTime(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func CoerceFromNullInt64(i sql.NullInt64) int64 {
	if i.Valid {
		return i.Int64
	}
	return 0
}
