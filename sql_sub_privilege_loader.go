package security

import (
	"context"
	"database/sql"
)

type SqlSubPrivilegeLoader struct {
	DB    *sql.DB
	Query string
}

func NewSubPrivilegeLoader(db *sql.DB, query string, options ...bool) *SqlSubPrivilegeLoader {
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
	return &SqlSubPrivilegeLoader{DB: db, Query: query}
}
func (l SqlSubPrivilegeLoader) Privilege(ctx context.Context, userId string, privilegeId string, sub string) int32 {
	var permissions int32 = 0
	err := l.DB.QueryRow(l.Query, userId, privilegeId, sub).Scan(&permissions)
	if err != nil {
		return ActionNone
	}
	if permissions == ActionNone {
		return ActionAll
	}
	return permissions
}
