package server

import (
	"github.com/gorilla/websocket"
	"log"
	"../fastJSON"
)

func appendVisitor(conn *websocket.Conn, msg string) {
	jsonMsg, err := parser.Parse(msg)
	if err != nil {
		log.Println("Error in visitor parse ", err)
		return
	}

	ok := jsonMsg.Exists("IMEI")
	if !ok {
		log.Println("Error in visitor parse IMEI don't exists")
		return
	}
	imei := jsonMsg.GetString("IMEI")

	ok = jsonMsg.Exists("table")
	if !ok {
		log.Println("Error in visitor parse table don't exists")
		return
	}
	table := jsonMsg.GetInt("table")
	if table == 0 {
		log.Println("Error in visitor parse table number don't number")
		return
	}

	visitors[conn] = visitorInfo{
		imei:imei,
		table:table,
	}
}

func isVisitor(jsonMsg *fastjson.Value) bool {
	ok := jsonMsg.Exists("type")
	if !ok {
		return false
	}

	if jsonMsg.GetString("type") == "client-authorization" {
		return true
	}

	return false
}

func visitorProcessing(conn *websocket.Conn) {

}