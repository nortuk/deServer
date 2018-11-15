package common

import (
	"github.com/gorilla/websocket"
)

var(
	DBConfig  DbCfg
	DBConnStr string

	ServConfig ServCfg

	StaffCon = make(map[*websocket.Conn]StaffInfo)
	VisitorsConn = make(map[*websocket.Conn]VisitorInfo)
	Tables = make(map[int]TableInfo)

	Menu = []MenuCategory{}
)
