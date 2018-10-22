package staff

import (
	"github.com/gorilla/websocket"
	"log"
	"svn.cloudserver.ru/fastJSON"
	"errors"
	"../common"
	"time"
)

func Processing(conn *websocket.Conn) error {
	for {
		msg, err := getMsg(conn)
		if err != nil {
			log.Println("Error in staff processing ", err)
			gapConn(conn)
			return err
		}

		command := msg.GetString("command")
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
	if (msgType == websocket.CloseMessage) || (err != nil) {
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

func gapConn(conn *websocket.Conn) {
	pers := common.StaffCon[conn]
	common.GapStaff[time.Now()] = pers
	delete(common.StaffCon, conn)
}