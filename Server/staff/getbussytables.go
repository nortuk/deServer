package staff

import (
	"github.com/gorilla/websocket"
	"log"
	"../common"
	"encoding/json"
)

func getBusyTables(conn *websocket.Conn) {
	var busyTables []int
	for id,info := range common.Tables {
		if len(info.Visitors) != 0 {
			busyTables = append(busyTables,id)
		}
	}

	if !sendBusyTables(conn, busyTables) {
		log.Println("Error in sending busy tables")
	}
}

func sendBusyTables(conn *websocket.Conn, busyTables []int) bool {
	answer := common.Response{
		Command: common.CommandGetBusyTables,
		Status: true,
		Data: common.DataStruct{
			"value": busyTables,
		},
	}

	jsonAnswer, err := json.Marshal(answer)
	if err != nil {
		log.Println("ERROR in marshal response:", err)
		return false
	}

	err = conn.WriteMessage(websocket.TextMessage,jsonAnswer)
	if err != nil {
		log.Println("ERROR in sending message:", err)
		return false
	}

	return true
}
