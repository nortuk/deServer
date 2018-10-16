package staff

import (
	"github.com/gorilla/websocket"
	"log"
	"database/sql"
	"encoding/json"
	"../common"
)

type table struct {
	TableID int `json:"id"`
	Name string `json:"name"`
}

func getTables(conn *websocket.Conn) {
	if len(common.Tables) == 0 {
		err := updateTables()
		if err != nil {
			log.Println("Error in gettables: ", err)
			common.SendError(conn, common.CommandGettables, common.ErrorDBProblem)
			return
		}
	}

	if !sendTables(conn) {
	log.Println("Error in sendtables")
	}
}

func sendTables(conn *websocket.Conn) bool{
	var tabs []table
	for id, tab := range common.Tables {
		tabs = append(tabs, table{
			TableID: id,
			Name: tab.Name,
		})
	}

	answer := common.Response{
		Command: "gettables",
		Status: true,
		Data: common.DataStruct{
			"value":tabs,
		},
	}

	jsonAnser, err := json.Marshal(answer)
	if err != nil {
		log.Println("ERROR in marshal response:", err)
		return false
	}

	err = conn.WriteMessage(websocket.TextMessage,jsonAnser)
	if err != nil {
		log.Println("ERROR in sending message:", err)
		return false
	}

	return true
}

func updateTables() error {
	sqlTables, err := getTablesFromDB()
	if err != nil {
		log.Println("Error in gettables: ", err)
		return err
	}
	for id,name := range sqlTables {
		common.Tables[id] = common.TableInfo{name, []string{}}
	}
	return nil
}

func getTablesFromDB() (tables map[int]string, err error) {
	db, err := sql.Open(common.DBConfig.UsedDatabase, common.DBConnStr)
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