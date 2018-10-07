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

	login := jsonMsg.GetString("login")
	if login == "" {
		log.Println("Error in visitor parse value don't exists")
		return
	}

	pass := jsonMsg.GetString("pass")
	if pass == "" {
		log.Println("Error in visitor parse password uncorrect")
		return
	}

	ok := database.CheckStaff(login, pass)
	if !ok {
		log.Println("Error uncorrect login or password")
		return
	}

	staff[conn] = staffInfo{
		login: login,
	}
}

func staffProcessing(conn *websocket.Conn) {
	for {
		return
	}
}