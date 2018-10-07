package database

import (
	"../Config"
	"log"
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"context"
)

var (
	cfg config.DBCfg
	connStr string
	ctx = context.Background()
)

func InitializeDBWork(path string) (er error) {
	var err error
	cfg, err = config.LoadDBCfg(path)
	if err != nil {
		log.Println("Error in initialize database work:", err)
		return err
	}

	connStr = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.NameDatabase)

	err = checkDBConnection()
	if err != nil {
		log.Println("Error in initialize database work:", err)
		return err
	}

	return nil
}


func checkDBConnection() (err error)  {
	db, err := sql.Open(cfg.UsedDatabase, connStr)
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("Error in ping database:", err)
		return err
	}

	return nil
}


func CheckStaff(login string, pass string) bool {
	db, err := sql.Open(cfg.UsedDatabase, connStr)
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		return false
	}
	defer db.Close()

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM staff WHERE login = $1 AND pass=$2`, login, pass).Scan(&count)
	if err != nil {
		log.Println("Error in request execution:", err)
		return false
	}

	if count == 1 {
		return true
	}

	return false
}