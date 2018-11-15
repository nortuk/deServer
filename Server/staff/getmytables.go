package staff

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"log"
	"../common"
)

func getMyTables(conn *websocket.Conn)  {
	common.StaffConnMutex.Lock()
	personal := common.StaffConn[conn]
	common.StaffConnMutex.Unlock()

	answer := common.Response{
		Command: "gettables",
		Status: true,
		Data: common.DataStruct{
			"value": personal.Tables,
		},
	}

	jsonAnser, err := json.Marshal(answer)
	if err != nil {
		log.Println("[" + conn.RemoteAddr().String() +"]ERROR in marshal response:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage,jsonAnser)
	if err != nil {
		log.Println("[" + conn.RemoteAddr().String() +"]ERROR in sending message:", err)
		return
	}

	return
}
