package staff

import (
	"github.com/gorilla/websocket"
	"../common"
	"log"
	"encoding/json"
)

func logout(conn *websocket.Conn) {
	personal := common.StaffCon[conn]

	for _, tableId := range personal.Tables {
		delete(common.Tables, tableId)
	}

	delete(common.StaffCon, conn)

	if !sendLogoutOK(conn) {
		log.Println("Error in sending ok command(settables)")
	}
}

func sendLogoutOK(conn *websocket.Conn) bool {
	answer := common.Response{
		Command: "logout",
		Status: true,
		Data: nil,
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
