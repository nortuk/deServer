package staff

import (
	"github.com/gorilla/websocket"
	"log"
	"../../fastJSON"
	"errors"
	"../common"
)

func Processing(msg *fastjson.Value,conn *websocket.Conn) error {
	err := auth(msg, conn)
	if err != nil {
		log.Println("[" + conn.RemoteAddr().String() + "]Error in authentification: ", err)
		return err
	}

	for {
		msg, err := getMsg(conn)
		if err != nil {
			return err
		}

		command := msg.GetString("command")

		_, ok := common.StaffCon[conn]
		if !ok {
			log.Println("[" + conn.RemoteAddr().String() +"]Connection close!")
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

		case common.CommandGetmenu:
			getMenu(conn)

		case common.CommandSetmenu:
			setmenu(msg,conn)

		case common.CommandGetTableInfo:
			gettableinfo(msg, conn)

		case common.CommandLogout:
			logout(conn)
			log.Println("[" + conn.RemoteAddr().String() +"]Logout")
			return nil

		default:
			common.SendError(conn, command, common.ErrorUnknownCommandType)
			log.Println("[" + conn.RemoteAddr().String() +"]Error: accept message with wrong structure")
		}
	}
}

func getMsg(conn *websocket.Conn) (msg *fastjson.Value, err error) {
	var parser = fastjson.Parser{}
	msgType, msgBytes, err := conn.ReadMessage()
	log.Println("[" + conn.RemoteAddr().String() +"]Accept message:", string(msgBytes))
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