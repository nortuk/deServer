package server

import (
	"github.com/gorilla/websocket"
	"log"
	"../Database"
	"svn.cloudserver.ru/fastJSON"
	"errors"
)

func appendStaff(conn *websocket.Conn, msg string) {
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
		case "getTables":
			getTables(conn)

		case "setTables":
			setTables(conn, msg)

		default:
			log.Println("Error: accept message with wrong structure")
		}
	}
}

func getTables(conn *websocket.Conn) {
	if len(tables) == 0 {
		err := updateTables()
		if err != nil {
			log.Println("Error in getTables: ", err)
			sendError(conn,"getTables", err.Error())
			return
		}
	}

	if !sendTables(conn) {
		log.Println("Error in sendTables")
		sendError(conn,"getTables","Error in sendTables")
	}
}

func updateTables() error {
	sqlTables, err := database.GetTables()
	if err != nil {
		log.Println("Error in getTables: ", err)
		return err
	}
	for id,name := range sqlTables {
		tables[id] = tableInfo{name, }
	}
	return nil
}

func setTables(conn *websocket.Conn, msg *fastjson.Value) {
	if len(tables) == 0 {
		err := updateTables()
		if err != nil {
			log.Println("Error in getTables: ", err)
			sendError(conn,"setTables", err.Error())
			return
		}
	}

	tableNumbers,err := getTablesNumbers(msg)
	if err != nil {
		log.Println("Error in getTablesNumbers: ", err)
		sendError(conn,"setTables", err.Error())
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
			sendError(conn,"setTables", "Number of table don't exists")
			return
		}
		personal.tables = append(personal.tables,num)
	}

	if !sendSetTablesOK(conn) {
		log.Println("Error in sending ok command(setTables)")
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