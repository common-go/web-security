package security

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

const (
	DriverPostgres   = "postgres"
	DriverMysql      = "mysql"
	DriverMssql      = "mssql"
	DriverOracle     = "oracle"
	DriverNotSupport = "no support"
)

type SqlPrivilegeLoader struct {
	DB    *sql.DB
	Query string
}

func NewPrivilegeLoader(db *sql.DB, query string) *SqlPrivilegeLoader {
	driver := GetDriver(db)
	query = replaceQueryparam(driver, query)
	return NewSqlPrivilegeLoader(db, query, driver)
}

func NewSqlPrivilegeLoader(db *sql.DB, query string, driver string) *SqlPrivilegeLoader {
	query = replaceQueryparam(driver, query)
	return &SqlPrivilegeLoader{DB: db, Query: query}
}

func (l SqlPrivilegeLoader) Privilege(ctx context.Context, userId string, privilegeId string) int32 {
	var permissions int32 = 0
	err := l.DB.QueryRow(l.Query, userId, privilegeId).Scan(&permissions)
	if err != nil {
		return ActionNone
	}
	if permissions == ActionNone {
		return ActionAll
	}
	return permissions
}

func replaceQueryparam(driver string, query string) string {
	if driver == DriverOracle || driver == DriverPostgres {
		var x string
		if driver == DriverOracle {
			x = ":val"
		} else {
			x = "$"
		}
		for i := 1; i < 3; i++ {
			query = strings.Replace(query, "?", x+fmt.Sprintf("%v", i), 1)
		}
	}
	return query
}

func GetDriver(db *sql.DB) string {
	driver := reflect.TypeOf(db.Driver()).String()
	switch driver {
	case "*pq.Driver":
		return DriverPostgres
	case "*mysql.MySQLDriver":
		return DriverMysql
	case "*mssql.Driver":
		return DriverMssql
	case "*godror.drv":
		return DriverOracle
	default:
		return DriverNotSupport
	}
}
