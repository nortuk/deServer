package server

import (
	"encoding/json"
	"log"
	"os"
	"./common"
)

func loadServCfg(path string) (cfg common.ServCfg, err error) {
	configFile, err := os.Open(path)
	if err != nil {
		log.Println("Error in the open of config server file:", err)
		return cfg, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&cfg)
	return cfg, nil
}

func loadDBCfg(path string) (cfg common.DbCfg, err error) {
	configFile, err := os.Open(path)
	if err != nil {
		log.Println("Error in the open of config database file:", err)
		return cfg, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&cfg)
	return cfg, nil
}
