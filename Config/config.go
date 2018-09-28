package config

import (
	"os"
	"log"
	"encoding/json"
)

type (
	DBCfg struct {
		UsedDatabase 	string `json:"usedDatabase"`
		NameDatabase 	string `json:"nameDatabase"`
		User 			string `json:"user"`
		Password 		string `json:"password"`
		Host 			string `json:"host"`
		Port 			string `json:"port"`
	}

	ServCfg struct {
		Ip 				string `json:"ip"`
		Port 			string `json:"port"`
	}
)

func LoadServCfg(path string) (cfg ServCfg, err error) {
	configFile, err := os.Open(path)
	if err != nil {
		log.Println("Error in the open of config server file:", err)
		return cfg, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&cfg)
	return cfg, nil
}

func LoadDBCfg(path string) (cfg DBCfg, err error) {
	configFile, err := os.Open(path)
	if err != nil {
		log.Println("Error in the open of config database file:", err)
		return cfg, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&cfg)
	return cfg,nil
}