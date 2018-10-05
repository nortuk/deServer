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

	connStr = fmt.Sprint("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable)",
		cfg.User, cfg.Password, cfg.NameDatabase, cfg.Host, cfg.Port)

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

	/*
	err = db.Ping()
	if err != nil {
		log.Println("Error in ping database:", err)
		return err
	}
	*/

	return nil
}


func CheckStaff(login string, pass string) bool {

	return true
	/*
	db, err := sql.Open(cfg.UsedDatabase, connStr)
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		return false
	}
	defer db.Close()

	//rows, err := db.QueryContext(ctx, "SELECT * FROM staff WHERE login=")
	if err != nil {
		log.Println("Error in request execution", err)
		return false
	}
	defer rows.Close()

*/
}