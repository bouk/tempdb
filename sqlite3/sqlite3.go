package sqlite3

import (
	"database/sql"
	"fmt"
	"sync/atomic"

	_ "github.com/mattn/go-sqlite3"
)

var nextID uint64 = 1

func New() (*sql.DB, func(), error) {
	id := atomic.AddUint64(&nextID, 1)
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%d?mode=memory&cache=shared", id))
	if err != nil {
		return nil, nil, err
	}

	return db, func() {
		db.Close()
	}, nil
}
