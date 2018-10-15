package common

import "github.com/gorilla/websocket"

var(
	DBConfig  DbCfg
	DBConnStr string

	ServConfig ServCfg

	StaffCon = make(map[*websocket.Conn]StaffInfo)
	Tables = make(map[int]TableInfo)
)
