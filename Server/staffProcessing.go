package server

import (
	"github.com/gorilla/websocket"
	"log"
	"../Database"
)

func appendStaff(conn *websocket.Conn, msg string) {
	jsonMsg, err := parser.Parse(msg)
	if err != nil {
		log.Println("Error in visitor parse ", err)
		return
	}

	ok := jsonMsg.Exists("data")
	if !ok {
		log.Println("Error in visitor parse value don't exists")
		return
	}
	data := jsonMsg.Get("data")

	login := data.GetString("login")
	if login == "" {
		log.Println("Error in visitor parse value don't exists")
		return
	}

	pass := data.GetString("pass")
	if pass == "" {
		log.Println("Error in visitor parse password uncorrect")
		return
	}

	id, err := database.GetStaffId(login, pass)
	if err != nil {
		log.Println("Error uncorrect login or password ", err)
		return
	}

	staff[conn] = staffInfo{
		id: id,
		login: login,
	}
}

func staffProcessing(conn *websocket.Conn) {
	for {
		msg, err := getMsg(conn)
		if err != nil {
			log.Println("Error in staff processing ", err)
			return
		}

		command := msg.GetString("command")
		switch command {
		case "getTables":
			getTables(conn)

		default:
			log.Println("Error: accept message with wrong structure")
		}
	}
}

func getTables(conn *websocket.Conn) {

}