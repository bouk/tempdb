package mysql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func New() (*sql.DB, func(), error) {
	dir, err := ioutil.TempDir("", "temporary-mysql")
	if err != nil {
		return nil, nil, err
	}
	cmd := exec.Command("mysqld", "--initialize-insecure", "-h", dir, "--skip-innodb-flush-sync", "--sync-binlog=0", "--sync-master-info=0", "--sync-relay-log=0", "--innodb-flush-method=nosync")
	cmd.Stderr = os.Stderr
	cmd.Stdout = ioutil.Discard
	if err = cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return nil, nil, err
	}

	socket := filepath.Join(dir, "mysql.sock")

	cmd = exec.Command("mysqld", "--skip-networking", "-h", dir, "--socket="+socket, "--pid-file="+dir+"/mysql.pid", "--mysqlx=OFF", "--skip-innodb-flush-sync", "--sync-binlog=0", "--sync-master-info=0", "--sync-relay-log=0", "--innodb-flush-method=nosync")
	cmd.Stderr = os.Stderr
	cmd.Stdout = ioutil.Discard
	if err = cmd.Start(); err != nil {
		os.RemoveAll(dir)
		return nil, nil, err
	}

	for i := 0; i < 100; i++ {
		if _, err := os.Stat(socket); err == nil {
			break
		} else if !os.IsNotExist(err) {
			cmd.Process.Kill()
			cmd.Wait()
			os.RemoveAll(dir)
			return nil, nil, err
		}
		time.Sleep(time.Millisecond * 25)
	}

	temp, err := sql.Open("mysql", fmt.Sprintf("root@unix(%s)/", socket))
	if err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		os.RemoveAll(dir)
		return nil, nil, err
	}

	_, err = temp.Exec("CREATE DATABASE data")
	if err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		os.RemoveAll(dir)
		return nil, nil, err
	}
	temp.Close()

	db, err := sql.Open("mysql", fmt.Sprintf("root@unix(%s)/data", socket))
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
