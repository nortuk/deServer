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
		tab := common.Tables[tableId]
		delete(tab.Staff, personal.Id)
		common.Tables[tableId] = tab
	}

	delete(common.StaffCon, conn)

	if !sendLogoutOK(conn) {
		log.Println("[" + conn.RemoteAddr().String() +"]Error in sending ok command(settables)")
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
		log.Println("[" + conn.RemoteAddr().String() +"]ERROR in marshal response:", err)
		return false
	}

	err = conn.WriteMessage(websocket.TextMessage,jsonAnser)
	if err != nil {
		log.Println("[" + conn.RemoteAddr().String() +"]ERROR in sending message:", err)
		return false
	}

	return true
}
