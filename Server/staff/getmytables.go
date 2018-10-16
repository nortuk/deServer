package staff

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"log"
	"../common"
)

func getMyTables(conn *websocket.Conn)  {
	personal := common.StaffCon[conn]

	answer := common.Response{
		Command: "gettables",
		Status: true,
		Data: common.DataStruct{
			"value": personal.Tables,
		},
	}

	jsonAnser, err := json.Marshal(answer)
	if err != nil {
		log.Println("ERROR in marshal response:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage,jsonAnser)
	if err != nil {
		log.Println("ERROR in sending message:", err)
		return
	}

	return
}
