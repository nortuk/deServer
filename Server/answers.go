package server

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"log"
)

type (
	dataStruct map[string]interface{}

	response struct {
		Command string `json:"command"`
		Status bool `json:"status"`
		Data dataStruct `json:"data"`
	}

	table struct {
		TableID int `json:"id"`
		Name string `json:"name"`
	}
)

func sendError(conn *websocket.Conn, command string, errMsg string) bool {
	answer := response{
		Command: command,
		Status: false,
		Data: dataStruct{
			"value": errMsg,
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

func sendAuthOk(conn *websocket.Conn) bool {
	answer := response{
		Command: "auth",
		Status: true,
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

func sendTables(conn *websocket.Conn) bool{
	var tabs = []table {}
	for id, tab := range tables {
		tabs = append(tabs, table{
			TableID: id,
			Name: tab.name,
		})
	}

	answer := response{
		Command: "gettables",
		Status: true,
		Data: dataStruct{
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

func sendSetTablesOK(conn *websocket.Conn) bool {
	answer := response{
		Command: "settables",
		Status: true,
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