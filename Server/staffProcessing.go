package server

import (
	"github.com/gorilla/websocket"
	"log"
)

func appendStaff(conn *websocket.Conn, msg string) {
	jsonMsg, err := parser.Parse(msg)
	if err != nil {
		log.Println("Error in visitor parse ", err)
		return
	}

	ok := jsonMsg.Exists("value")
	if !ok {
		log.Println("Error in visitor parse value don't exists")
		return
	}
	value := jsonMsg.Get("value")

	login := value.GetString("login")
	if login == "" {
		log.Println("Error in visitor parse login uncorrect")
		return
	}

	pass := value.GetString("password")
	if pass == "" {
		log.Println("Error in visitor parse password uncorrect")
		return
	}

	ok = database.CheckStaff(login, pass)
	if !ok {
		log.Println("Error uncorrect login or password")
		return
	}

	staff[conn] = staffInfo{
		login: login,
	}
}