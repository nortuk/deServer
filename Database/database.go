package database

import (
	"../Config"
	"log"
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"context"
	"errors"
)

var (
	Cfg     config.DBCfg
	ConnStr string
	ctx   = context.Background()
)

func InitializeDBWork(path string) (er error) {
	var err error
	Cfg, err = config.LoadDBCfg(path)
	if err != nil {
		log.Println("Error in initialize database work:", err)
		return err
	}

	ConnStr = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		Cfg.Host, Cfg.Port, Cfg.User, Cfg.NameDatabase)

	err = checkDBConnection()
	if err != nil {
		log.Println("Error in initialize database work:", err)
		return err
	}

	return nil
}


func checkDBConnection() (err error)  {
	db, err := sql.Open(Cfg.UsedDatabase, ConnStr)
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


func GetStaffId(login string, pass string) (id int, err error) {
	db, err := sql.Open(Cfg.UsedDatabase, ConnStr)
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		return 0, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT ID FROM staff WHERE login = $1 AND pass=$2`, login, pass)
	if err != nil {
		log.Println("Error in request execution:", err)
		return 0,err
	}
	var userID sql.NullInt64
	if rows.Next() {
		rows.Scan(&userID)
		return int(userID.Int64), nil
	}
	return 0, errors.New("Error: can't find this user")
}

func GetTables() (tables map[int]string, err error) {
	db, err := sql.Open(Cfg.UsedDatabase, ConnStr)
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		return nil,err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, name FROM tables`)
	if err != nil {
		log.Println("Error in request execution:", err)
		return nil,err
	}
	var id sql.NullInt64
	var name sql.NullString
	res := make(map[int]string)
	for rows.Next() {
		err = rows.Scan(&id,&name)
		if err != nil {
			log.Println("Error: in reading tables: ", err)
			return nil, err
		}
		res[int(id.Int64)] = name.String
	}

	return res, nil
}

