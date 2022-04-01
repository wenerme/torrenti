package nilx

import (
	"database/sql"
	"time"
)

func StringSQLNull(v sql.NullString) *string {
	if v.Valid {
		return &v.String
	}
	return nil
}

func TimeSQLNull(v sql.NullTime) *time.Time {
	if v.Valid {
		return &v.Time
	}
	return nil
}

func EmptyToSQLNull(s string) (o sql.NullString) {
	o.String = s
	o.Valid = s != ""
	return
}
