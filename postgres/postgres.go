package postgres

import (
	"database/sql"
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
	cmd := exec.Command("initdb", "--nosync", dir)
	cmd.Stderr = os.Stderr
	cmd.Stdout = ioutil.Discard
	if err = cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return nil, nil, err
	}
	cmd = exec.Command("postgres", "-F", "--unix_socket_directories="+dir, "--listen_addresses=\"\"", "-D", dir)
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
		return nil, nil, err
	}
	return db, func() {
		db.Close()
		cmd.Process.Kill()
		cmd.Wait()
		os.RemoveAll(dir)
		return
	}, nil
}
