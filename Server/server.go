package server

import (
	"../fastJSON"
	"errors"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"./staff"
	"./common"
	"./visitor"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true //чтоб разрешить кроссдоменные запросы
		},
	}
)

func Start(servCfgPath string, dbCfgPath string) (err error) {
	err = initializeDBWork(dbCfgPath)
	if err != nil {
		log.Println("EROR in server start: database initialization: ", err)
		return err
	}

	common.ServConfig, err = loadServCfg(servCfgPath)
	if err != nil {
		log.Println("ERROR in loading server config: ", err)
		return err
	}

	listening()

	return nil
}

func listening() {
	err := loadTables()
	if err != nil {
		log.Fatal("ERROR in loading tables from DB", err)
	}
	err = loadMenu()
	if err != nil {
		log.Fatal("ERROR in loading menu from DB", err)
	}

	log.Println("Server started (" + common.ServConfig.Ip + ":" + common.ServConfig.Port + ")")
	http.HandleFunc("/", handler)
	addr := flag.String("addr",
		common.ServConfig.Ip+":"+ common.ServConfig.Port, "localhost")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	//из http запроса переходим в websocket соединение
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error in upgrading http response from " + r.RemoteAddr + " : ", err)
		return
	}

	listeningConnection(conn)
}

func listeningConnection(conn *websocket.Conn) {
	for {
		msg, err := getMsg(conn)
		if err != nil {
			log.Println("[" + conn.RemoteAddr().String() + "]Error in geting message: ", err)
			common.SendError(conn, "", common.ErrorUnknownMsg)
			return
		}

		conType := getAuthType(msg)

		switch conType {
		case isVisitorType:
			err := visitor.Processing(msg, conn)
			if err != nil {
				log.Println("[" + conn.RemoteAddr().String() +"]Error in visitor processing: ", err)
				return
			}

		case isStaffType:
			err := staff.Processing(msg, conn)
			if err != nil{
				log.Println("[" + conn.RemoteAddr().String() +"]Error in staff processing: ", err)
				return
			}

		case unknownType:
			common.SendError(conn, "", common.ErrorUnknownCommandType)
			log.Println("Unknown type")
		}

	}
}

func getMsg(conn *websocket.Conn) (msg *fastjson.Value, err error) {
	var parser = fastjson.Parser{}
	msgType, msgBytes, err := conn.ReadMessage()
	log.Println("[" + conn.RemoteAddr().String() + "]Accept message: ", string(msgBytes))
	if msgType == websocket.CloseMessage {
		log.Println("----[" + conn.RemoteAddr().String() + "]WebSocket close message" )
		msg, _ := parser.Parse("{}")
		return msg, errors.New("Connection closed")
	}
	if err != nil {
		log.Println("----[" + conn.RemoteAddr().String() +  "]Error in get msg:", err)
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