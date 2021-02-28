package security

import (
	"context"
	"database/sql"
	"fmt"
)

type SqlPrivilegesLoader struct {
	DB    *sql.DB
	Query string
}

func NewPrivilegesLoader(db *sql.DB, query string, options ...bool) *SqlPrivilegesLoader {
	var handleDriver bool
	if len(options) >= 1 {
		handleDriver = options[0]
	} else {
		handleDriver = true
	}
	if handleDriver {
		driver := GetDriver(db)
		query = ReplaceQueryArgs(driver, query)
	}
	return &SqlPrivilegesLoader{DB: db, Query: query}
}

func (l SqlPrivilegesLoader) Privileges(ctx context.Context, userId string) []string {
	privileges := make([]string, 0)
	rows, err := l.DB.Query(l.Query, userId)
	if err != nil {
		return privileges
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var permissions int32
		if err = rows.Scan(&id, &permissions); err == nil {
			if permissions != ActionNone {
				x := id + " " + fmt.Sprintf("%X", permissions)
				privileges = append(privileges, x)
			} else {
				privileges = append(privileges, id)
			}
		}
	}
	return privileges
}
