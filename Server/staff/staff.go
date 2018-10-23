package staff

import (
	"github.com/gorilla/websocket"
	"log"
	"svn.cloudserver.ru/fastJSON"
	"errors"
	"../common"
)

func Processing(conn *websocket.Conn) error {
	for {
		msg, err := getMsg(conn)
		if err != nil {
			log.Println("Error in staff processing ", err)
			return err
		}

		command := msg.GetString("command")

		_, ok := common.StaffCon[conn]
		if !ok {
			log.Println("Connection close!")
			common.SendError(conn, command, common.ErrorConnectionrefused)
			return nil
		}

		switch command {
		case common.CommandGettables:
			getTables(conn)

		case common.CommandSettables:
			setTables(conn, msg)

		case common.CommandGetMyTables:
			getMyTables(conn)

		case common.CommandGetBusyTables:
			getBusyTables(conn)

		case common.CommandLogout:
			logout(conn)
			return nil

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
	if (msgType != websocket.TextMessage) || (err != nil) {
		msg, _ := parser.Parse("{}")
		return msg, errors.New("Connection closed")
	}

	msg, parseErr := parser.Parse(string(msgBytes))
	if parseErr != nil {
		msg, _ := parser.Parse("{}")
		return msg, nil
	}

	return msg, nil
}