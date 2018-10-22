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
	"database/sql"
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
		log.Println("EROR(Server start): database initialization: ", err)
		return err
	}

	common.ServConfig, err = loadServCfg(servCfgPath)
	if err != nil {
		return err
	}

	listening()

	return nil
}

func listening() {
	err := loadTables()
	if err != nil {
		log.Fatal("Can't load tables from DB")
	}

	log.Println("Server started")
	addr := flag.String("addr",
		common.ServConfig.Ip+":"+ common.ServConfig.Port, "localhost")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*addr, nil))
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
		msg, err := getMsg(conn)
		if err != nil {
			log.Println("Error in geting message: ", err)
			common.SendError(conn, "", common.ErrorUnknownMsg)
			return
		}

		conType := getAuthType(msg)

		switch conType {
		case isVisitorType:
			return

		case isStaffType:
			if staffAuth(msg, conn) {
				err := staff.Processing(conn)
				if err != nil{
					log.Println("Error in staff processing: ", err)
					return
				}
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

func loadTables() error {
	db, err := sql.Open(common.DBConfig.UsedDatabase, common.DBConnStr)
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		return err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, name FROM tables`)
	if err != nil {
		log.Println("Error in request execution:", err)
		return err
	}
	var id sql.NullInt64
	var name sql.NullString
	for rows.Next() {
		err = rows.Scan(&id,&name)
		if err != nil {
			log.Println("Error: in reading tables: ", err)
			return err
		}
		common.Tables[int(id.Int64)] = common.TableInfo{
			Name: name.String,
			Visitors: []string{},
			Staff: map[int]int{},
		}
	}

	return nil
}