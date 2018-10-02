package server

import (
	"../Database"
	"log"
	"../Config"
	"github.com/gorilla/websocket"
	"net/http"
	"flag"
	"svn.cloudserver.ru/fastJSON"
)

type (
	staffInfo struct {
		login string
	}

	visitorInfo struct {
		imei string
		table int
	}

	personType int
)

const (
	isVisitorType personType = 1
	isStaffType personType = 2
	unknownType personType = 3
)

var (
	parser = fastjson.Parser{}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true //чтоб разрешить кроссдоменные запросы
		},
	}

	cfg config.ServCfg

	staff = make(map[*websocket.Conn]staffInfo)
	visitors = make(map[*websocket.Conn]visitorInfo)

)

func Start(servCfgPath string, dbCfgPath string) (err error) {
	err = database.InitializeDBWork(dbCfgPath)
	if err != nil {
		log.Println("Error in database initialization:", err)
		return err
	}
	
	cfg, err = config.LoadServCfg(servCfgPath)
	if err != nil {
		return err
	}
	
	listening()
	
	return nil
}

func listening()  {
	http.HandleFunc("/", handler)
	log.Println("Server started")
	value := cfg.Ip + ":" + cfg.Port
	addr := flag.String("addr", value, "localhost")
	log.Fatal(http.ListenAndServe(*addr,nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error", err)
		return
	}

	listeningConnection(conn)
}

func listeningConnection(conn *websocket.Conn) {
	for {
		msgType, msg, err := conn.ReadMessage()
		if !checkCorrectMsg(msgType, err) {
			return
		}

		personalType := getPersonalType(string(msg))
		switch personalType {
		case isVisitorType:
			appendVisitor(conn,string(msg))
			_, ok := visitors[conn]
			if !ok {
				log.Println("Error in parsing info")
				continue
			}

			visitorProcessing(conn)
			return

		case isStaffType:
			ok := checkPersonal(string(msg))
			if !ok {
				log.Println("Error in check personal")
				continue
			}


		case unknownType:
			log.Println("Unknown type")
		}
	}
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

func isVisitor(jsonMsg *fastjson.Value) bool {
	ok := jsonMsg.Exists("type")
	if !ok {
		return false
	}

	if jsonMsg.GetString("type") == "client-authorization" {
		return true
	}

	return false
}

func  isStaff(jsonMsg *fastjson.Value) bool {
	ok := jsonMsg.Exists("type")
	if !ok {
		return false
	}

	if jsonMsg.GetString("type") == "auth" {
		return true
	}

	return false
}

func appendVisitor(conn *websocket.Conn, msg string) {
	jsonMsg, err := parser.Parse(msg)
	if err != nil {
		log.Println("Error in visitor parse ", err)
		return
	}

	ok := jsonMsg.Exists("IMEI")
	if !ok {
		log.Println("Error in visitor parse IMEI don't exists")
		return
	}
	imei := jsonMsg.GetString("IMEI")

	ok = jsonMsg.Exists("table")
	if !ok {
		log.Println("Error in visitor parse table don't exists")
		return
	}
	table := jsonMsg.GetInt("table")
	if table == 0 {
		log.Println("Error in visitor parse table number don't number")
		return
	}

	visitors[conn] = visitorInfo{
		imei:imei,
		table:table,
	}
}

func checkPersonal(msg string) bool {
	return false
}

func visitorProcessing(conn *websocket.Conn) {

}