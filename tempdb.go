// Package tempdb allows you to quickly create fresh throw-away databases for testing
package tempdb // import "bou.ke/tempdb"

import (
	"database/sql"
	"fmt"

	"bou.ke/tempdb/mysql"
	"bou.ke/tempdb/postgres"
	"bou.ke/tempdb/sqlite3"
)

func New(driver string) (*sql.DB, func(), error) {
	switch driver {
	case "mysql":
		return mysql.New()
	case "postgres":
		return postgres.New()
	case "sqlite3":
		return sqlite3.New()
	default:
		return nil, nil, fmt.Errorf("unsupported driver %q", driver)
	}
}
