package staff

import (
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"../common"
)

type table struct {
	TableID int `json:"id"`
	Name string `json:"name"`
}

func getTables(conn *websocket.Conn) {
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