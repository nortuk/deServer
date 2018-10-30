package staff

import (
	"github.com/gorilla/websocket"
	"../../fastJSON"
	"log"
	"../common"
	"errors"
	"encoding/json"
	"database/sql"
	"strconv"
)

func setTables(conn *websocket.Conn, msg *fastjson.Value) {
	tableNumbers,err := getTablesNumbers(msg)
	if err != nil {
		log.Println("[" + conn.RemoteAddr().String() +"]Error in gettablesNumbers: ", err)
		common.SendError(conn,common.CommandSettables, common.ErrorWrongCommandStructure)
		return
	}

	personal := common.StaffCon[conn]
	oldTables := personal.Tables
	if len(personal.Tables) != 0 {
		personal.Tables = personal.Tables[:0]
	}
	for _,num := range tableNumbers {
		_, ok := common.Tables[num]
		if !ok {
			log.Println("[" + conn.RemoteAddr().String() +"]Error: number of table doesn't exists")
			common.SendError(conn,common.CommandSettables, common.ErrorTableDoesnotExists)
			personal.Tables = oldTables
			return
		}
		personal.Tables = append(personal.Tables,num)
	}

	if !setTablesInDB(personal) {
		log.Println("[" + conn.RemoteAddr().String() +"]Error in append table in DB")
		common.SendError(conn,common.CommandSettables, common.ErrorDBProblem)
		personal.Tables = oldTables
		return
	}

	for _, tableId := range oldTables {
		tab := common.Tables[tableId]
		delete(tab.Staff, personal.Id)
		common.Tables[tableId]= tab
	}

	for _, tableID := range personal.Tables {
		common.Tables[tableID].Staff[personal.Id] = personal.Id
	}


	if !sendSetTablesOK(conn) {
		log.Println("[" + conn.RemoteAddr().String() +"]Error in sending ok command(settables)")
	}
}

func sendSetTablesOK(conn *websocket.Conn) bool {
	answer := common.Response{
		Command: common.CommandSettables,
		Status: true,
		Data: nil,
	}

	jsonAnser, err := json.Marshal(answer)
	if err != nil {
		log.Println("[" + conn.RemoteAddr().String() +"]ERROR in marshal response:", err)
		return false
	}

	err = conn.WriteMessage(websocket.TextMessage,jsonAnser)
	if err != nil {
		log.Println("[" + conn.RemoteAddr().String() +"]ERROR in sending message:", err)
		return false
	}

	return true
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

func setTablesInDB(personal common.StaffInfo) bool {
	query := "DELETE FROM staff_tables WHERE id_staff = " +
		strconv.Itoa(personal.Id) + ";"
	for _, tableID := range personal.Tables {
		query += "INSERT INTO staff_tables(id_staff, id_table) VALUES (" +
			strconv.Itoa(personal.Id) + "," + strconv.Itoa(tableID) + ");"
	}

	db, err := sql.Open(common.DBConfig.UsedDatabase, common.DBConnStr)
	if err != nil {
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return false
	}

	_, err = stmt.Exec()
	if err != nil {
		return false
	}

	return true
}