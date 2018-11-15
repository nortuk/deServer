package common

import (
	"github.com/gorilla/websocket"
	"sync"
)

var (
	DBConfig  DbCfg
	DBConnStr string

	ServConfig ServCfg

	StaffConnMutex    = &sync.Mutex{}
	StaffConn         = make(map[*websocket.Conn]StaffInfo)
	VisitorsConnMutex = &sync.Mutex{}
	VisitorsConn      = make(map[*websocket.Conn]VisitorInfo)
	TablesMutex       = &sync.Mutex{}
	Tables            = make(map[int]TableInfo)

	Menu = []MenuCategory{}
)
