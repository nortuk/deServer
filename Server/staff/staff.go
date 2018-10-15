package staff

import (
	"github.com/gorilla/websocket"
	"log"
	"svn.cloudserver.ru/fastJSON"
	"errors"
	"../common"
)

func Processing(conn *websocket.Conn){
	for {
		msg, err := getMsg(conn)
		if err != nil {
			log.Println("Error in staff processing ", err)
			return
		}

		command := msg.GetString("command")
		switch command {
		case common.CommandGettables:
			getTables(conn)

		case common.CommandSettables:
			//setTables(conn, msg)

		case "getbusytables":
			//getBusyTables(conn)

		default:
			common.SendError(conn, command, common.ErrorUnknownCommandType)
			log.Println("Error: accept message with wrong structure")
		}
	}
}

func getMsg(conn *websocket.Conn) (msg *fastjson.Value, err error) {
	var parser = fastjson.Parser{}
	msgType, msgBytes, err := conn.ReadMessage()
	log.Println("Accept message:", string(msgBytes))
	if msgType == websocket.CloseMessage {
		//вызов функции при разрыве
		log.Println("Close message")
		msg, _ := parser.Parse("{}")
		return msg, errors.New("Connection closed")
	}
	if err != nil {
		//вызов функции при разрыве
		log.Println("Error in get msg:", err)
		msg, _ := parser.Parse("{}")
		return msg, err
	}

	msg, parseErr := parser.Parse(string(msgBytes))
	if parseErr != nil {
		msg, _ := parser.Parse("{}")
		return msg, nil
	}

	return msg, nil
}