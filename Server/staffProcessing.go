package server

import (
	"github.com/gorilla/websocket"
	"log"
	"../Database"
	"svn.cloudserver.ru/fastJSON"
	"errors"
)

func appendStaff(conn *websocket.Conn, msg string) {
	var parser = fastjson.Parser{}
	jsonMsg, err := parser.Parse(msg)
	if err != nil {
		log.Println("Error in visitor parse ", err)
		return
	}

	ok := jsonMsg.Exists("data")
	if !ok {
		log.Println("Error in visitor parse value don't exists")
		return
	}
	data := jsonMsg.Get("data")

	login := data.GetString("login")
	if login == "" {
		log.Println("Error in visitor parse value don't exists")
		return
	}

	pass := data.GetString("pass")
	if pass == "" {
		log.Println("Error in visitor parse password uncorrect")
		return
	}

	id, err := database.GetStaffId(login, pass)
	if err != nil {
		log.Println("Error uncorrect login or password ", err)
		return
	}

	staff[conn] = staffInfo{
		id: id,
		login: login,
	}
}

func staffProcessing(conn *websocket.Conn) {
	for {
		msg, err := getMsg(conn)
		if err != nil {
			log.Println("Error in staff processing ", err)
			return
		}

		command := msg.GetString("command")
		switch command {
		case "gettables":
			getTables(conn)

		case "settables":
			setTables(conn, msg)

		case "getbusytables":
			getBusyTables(conn)

		default:
			log.Println("Error: accept message with wrong structure")
		}
	}
}

func getTables(conn *websocket.Conn) {
	if len(tables) == 0 {
		err := updateTables()
		if err != nil {
			log.Println("Error in gettables: ", err)
			sendError(conn,"gettables", err.Error())
			return
		}
	}

	if !sendTables(conn) {
		log.Println("Error in sendtables")
		sendError(conn,"gettables","Error in sendtables")
	}
}

func updateTables() error {
	sqlTables, err := database.GetTables()
	if err != nil {
		log.Println("Error in gettables: ", err)
		return err
	}
	for id,name := range sqlTables {
		tables[id] = tableInfo{name, nil }
	}
	return nil
}

func setTables(conn *websocket.Conn, msg *fastjson.Value) {
	if len(tables) == 0 {
		err := updateTables()
		if err != nil {
			log.Println("Error in gettables: ", err)
			sendError(conn,"settables", err.Error())
			return
		}
	}

	tableNumbers,err := getTablesNumbers(msg)
	if err != nil {
		log.Println("Error in gettablesNumbers: ", err)
		sendError(conn,"settables", err.Error())
		return
	}

	personal := staff[conn]
	if len(personal.tables) != 0 {
		personal.tables = personal.tables[:0]
	}
	for _,num := range tableNumbers {
		_, ok := tables[num]
		if !ok {
			log.Println("Error number of table don't exists")
			sendError(conn,"settables", "Number of table don't exists")
			return
		}
		personal.tables = append(personal.tables,num)
	}

	if !sendSetTablesOK(conn) {
		log.Println("Error in sending ok command(settables)")
	}
}

func getTablesNumbers(msg *fastjson.Value) (val []int, err error) {
	if !msg.Exists("data") {
		return nil, errors.New("Don't exist data")
	}
	data := msg.Get("data")
	if !data.Exists("value") {
		return nil, errors.New("Don't exist value")
	}
	value, err := data.Get("value").Array()
	if err != nil {
		return nil, err
	}
	for _, number := range value {
		num, err := number.Int()
		if err != nil {
			return nil, err
		}
		val = append(val,num)
	}

	return val,err
}

func getBusyTables(conn *websocket.Conn) {
	var busyTables []int
	for id,info := range tables {
		if len(info.visitors) != 0 {
			busyTables = append(busyTables,id)
		}
	}

	if !sendBusyTables(conn, busyTables) {
		log.Println("Error in sending busy tables")
	}
}