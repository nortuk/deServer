package Config

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"log"
	"encoding/json"
	"fmt"
)

type (
	Database struct {
		UsedDatabase 	string `json:"usedDatabase"`
		NameDatabase 	string `json:"nameDatabase"`
		User 			string `json:"user"`
		Password 		string `json:"password"`
		Host 			string `json:"host"`
		Port 			string `json:"port"`
	}

	Server struct {
		Ip 				string `json:"ip"`
		Port 			string `json:"port"`
	}
)

func ReadDatabaseConfig (path string) (db *sql.DB, err error) {
	var config Database
	configFile, err := os.Open(path)
	if err != nil {
		log.Println("Error in the open of config database file:", err)
		return db, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	var dataSource string
	dataSource = fmt.Sprint("user=%s password=%s dbname=&s host=&s port=%s sslmode=disable)",
		config.User, config.Password, config.NameDatabase, config.Host, config.Port)
	db, err = sql.Open(config.UsedDatabase, dataSource)
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		return db, err
	}
	return db,nil
}

func ReadServerConfig (path string) (serv Server, err error) {
	configFile, err := os.Open(path)
	if err != nil {
		log.Println("Error in the open of config server file:", err)
		return serv, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&serv)
	return serv, nil
}

