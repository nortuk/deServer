package staff

import (
	"github.com/gorilla/websocket"
	"../common"
)

func logout(conn *websocket.Conn) {
	personal := common.StaffCon[conn]

	for _, tableId := range personal.Tables {
		delete(common.Tables, tableId)
	}

	delete(common.StaffCon, conn)
}
