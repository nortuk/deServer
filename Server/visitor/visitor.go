package visitor

import (
	"../../fastJSON"
	"github.com/gorilla/websocket"
	"log"
	"errors"
	"../common"
)

func Processing(msg *fastjson.Value, conn *websocket.Conn) error {
	err := auth(msg, conn)
	if err != nil {
		log.Println("[" + conn.RemoteAddr().String() + "]Error in visitor authentification: ", err)
		return nil
	}

	for {
		msg, err := getMsg(conn)
		if err != nil {
			return err
		}

		command := msg.GetString("command")

		_, ok := common.VisitorsConn[conn]
		if !ok {
			log.Println("[" + conn.RemoteAddr().String() +"]Connection close!")
			common.SendError(conn, command, common.ErrorConnectionrefused)
			return nil
		}

		switch command {
		case common.CommandGetmenu:
			getMenu(conn)



		default:
			common.SendError(conn, command, common.ErrorUnknownCommandType)
			log.Println("[" + conn.RemoteAddr().String() +"]Error: accept message with wrong structure")
		}
	}

	return nil
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