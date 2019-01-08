package tempdb

import (
	"database/sql"
	"testing"

	"gopkg.in/stretchr/testify.v1/require"
)

func TestAll(t *testing.T) {
	for _, driver := range sql.Drivers() {
		t.Run(driver, func(t *testing.T) {
			db, cleanup, err := New(driver)
			require.NoError(t, err)
			defer cleanup()
			row := db.QueryRow("SELECT 123")
			var n int
			require.NoError(t, row.Scan(&n))
			require.Equal(t, 123, n)
		})
	}
}
