package postgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

func New() (*sql.DB, func(), error) {
	dir, err := ioutil.TempDir("", "temporary-postgres")
	if err != nil {
		return nil, nil, err
	}
	cmd := exec.Command("initdb", "--nosync", "--encoding=UNICODE", "--auth=trust", dir)
	cmd.Stderr = os.Stderr
	cmd.Stdout = ioutil.Discard
	if err = cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return nil, nil, err
	}

	cmd = exec.Command(
		"postgres", "-F", "-D", dir,
		"-c", "unix_socket_directories="+dir,
		"-c", "listen_addresses=",
		"-c", "shared_buffers=12MB",
		"-c", "fsync=off",
		"-c", "synchronous_commit=off",
		"-c", "full_page_writes=off",
		"-c", "track_activities=off",
		"-c", "track_counts=off",
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = ioutil.Discard
	if err = cmd.Start(); err != nil {
		os.RemoveAll(dir)
		return nil, nil, err
	}

	for i := 0; i < 100; i++ {
		if _, err := os.Stat(filepath.Join(dir, ".s.PGSQL.5432")); err == nil {
			break
		} else if !os.IsNotExist(err) {
			cmd.Process.Kill()
			cmd.Wait()
			os.RemoveAll(dir)
			return nil, nil, err
		}
		time.Sleep(time.Millisecond * 25)
	}

	db, err := sql.Open("postgres", "dbname=postgres host="+dir)
	if err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		os.RemoveAll(dir)
		return nil, nil, fmt.Errorf("opening database: %w", err)
	}
	for i := 0; i < 1000; i++ {
		var n int
		if db.QueryRow("SELECT 1").Scan(&n) == nil {
			break
		}
		time.Sleep(time.Millisecond * 25)
	}
	return db, func() {
		db.Close()
		cmd.Process.Kill()
		cmd.Wait()
		os.RemoveAll(dir)
		return
	}, nil
}
