package server

import (
	"log"
	"svn.cloudserver.ru/fastJSON"
	"github.com/gorilla/websocket"
)

func getPersonalType(msg string) (res personType) {
	jsonMsg, err := parser.Parse(msg)
	if err != nil {
		log.Println("Error in authentification parse ", err)
		return unknownType
	}

	ok := isVisitor(jsonMsg)
	if ok {
		return isVisitorType
	}

	ok = isStaff(jsonMsg)
	if ok {
		return isStaffType
	}

	return unknownType
}

func  isStaff(jsonMsg *fastjson.Value) bool {
	ok := jsonMsg.Exists("command")
	if !ok {
		return false
	}

	if jsonMsg.GetString("command") == "auth" {
		return true
	}

	return false
}

func checkCorrectMsg(msgType int, err error) bool {
	if msgType == websocket.CloseMessage {
		return false
	}

	if err != nil {
		log.Println("Error while read message")
		return false
	}

	return true
}