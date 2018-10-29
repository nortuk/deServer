package staff

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"log"
	"../common"
)

func getMenu(conn *websocket.Conn) {
	answer := common.Response{
		Command: "getmenu",
		Status: true,
		Data: common.DataStruct{
			"value": common.Menu,
		},
	}

	jsonAnswer, err := json.Marshal(answer)
	if err != nil {
		log.Println("[" + conn.RemoteAddr().String() +"]ERROR in marshal response:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage,jsonAnswer)
	if err != nil {
		log.Println("[" + conn.RemoteAddr().String() +"]ERROR in sending message:", err)
		return
	}

	return
}