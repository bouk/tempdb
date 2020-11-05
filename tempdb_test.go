package tempdb

import (
	"database/sql"
	"testing"

	"gopkg.in/stretchr/testify.v1/require"
)

func TestAll(t *testing.T) {
	for _, driver := range sql.Drivers() {
		t.Run(driver, func(t *testing.T) {
			db := TestDB(t, driver)
			row := db.QueryRow("SELECT 123")
			var n int
			require.NoError(t, row.Scan(&n))
			require.Equal(t, 123, n)
		})
	}
}
