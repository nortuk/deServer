package server

import (
	"database/sql"
	_"github.com/lib/pq"
	"log"
	"fmt"
	"./common"
)

func initializeDBWork(path string) (err error) {
	common.DBConfig, err = loadDBCfg(path)
	if err != nil {
		log.Println("EROR: load database config:", err)
		return err
	}

	if common.DBConfig.Password == "" {
		common.DBConnStr = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
			common.DBConfig.Host, common.DBConfig.Port, common.DBConfig.User, common.DBConfig.NameDatabase)
	} else{
		common.DBConnStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			common.DBConfig.Host, common.DBConfig.Port, common.DBConfig.User,
			common.DBConfig.Password, common.DBConfig.NameDatabase)
	}

	err = checkDBConnection()
	if err != nil {
		log.Println("Error in initialize database work:", err)
		return err
	}

	return nil
}

func checkDBConnection() (err error) {
	db, err := sql.Open(common.DBConfig.UsedDatabase, common.DBConnStr)
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
