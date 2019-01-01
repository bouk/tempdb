package tempdb

import (
	"database/sql"
	"fmt"

	"bou.ke/tempdb/postgres"
	"bou.ke/tempdb/sqlite3"
)

func New(driver string) (*sql.DB, func(), error) {
	switch driver {
	case "postgres":
		return postgres.New()
	case "sqlite3":
		return sqlite3.New()
	default:
		return nil, nil, fmt.Errorf("unsupported driver %q", driver)
	}
}
