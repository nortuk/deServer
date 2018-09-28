package database

import (
	"../Config"
	"log"
	"fmt"
	"database/sql"
)

var (
	cfg config.DBCfg
)

func InitializeDBWork(path string) (er error) {
	var err error
	cfg, err = config.LoadDBCfg(path)
	if err != nill {
		log.Print("Error in initialize database work:", err)
		return err
	}

	err = checkDBConnection()
	if err != nil {
		log.Print("Error in initialize database work:", err)
		return err
	}

	return nil
}

func checkDBConnection() (err error)  {
	connStr := fmt.Sprint("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable)",
		cfg.User, cfg.Password, cfg.NameDatabase, cfg.Host, cfg.Port)
	db, err := sql.Open(cfg.UsedDatabase, connStr)
	defer db.Close()
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		return err
	}

	return nil
}