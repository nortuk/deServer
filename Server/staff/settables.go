package staff

import (
	"github.com/gorilla/websocket"
	"svn.cloudserver.ru/fastJSON"
	"log"
	"../common"
	"errors"
	"encoding/json"
	"database/sql"
	"strconv"
)

func setTables(conn *websocket.Conn, msg *fastjson.Value) {
	if len(common.Tables) == 0 {
		err := updateTables()
		if err != nil {
			log.Println("Error in gettables: ", err)
			common.SendError(conn, common.CommandSettables, common.ErrorDBProblem)
			return
		}
	}

	tableNumbers,err := getTablesNumbers(msg)
	if err != nil {
		log.Println("Error in gettablesNumbers: ", err)
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
			log.Println("Error number of table don't exists")
			common.SendError(conn,common.CommandSettables, common.ErrorTableDoesnotExists)
			personal.Tables = oldTables
			return
		}
		personal.Tables = append(personal.Tables,num)
	}

	if !setTablesInDB(personal) {
		log.Println("Error in append table in DB")
		common.SendError(conn,common.CommandSettables, common.ErrorDBProblem)
		personal.Tables = oldTables
		return
	}

	if !sendSetTablesOK(conn) {
		log.Println("Error in sending ok command(settables)")
	}
}

func sendSetTablesOK(conn *websocket.Conn) bool {
	answer := common.Response{
		Command: "gettables",
		Status: true,
		Data: nil,
	}

	jsonAnser, err := json.Marshal(answer)
	if err != nil {
		log.Println("ERROR in marshal response:", err)
		return false
	}

	err = conn.WriteMessage(websocket.TextMessage,jsonAnser)
	if err != nil {
		log.Println("ERROR in sending message:", err)
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
		log.Println("Error in the open connection with database:", err)
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("Error in creating query:", err)
		return false
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Println("Error in executing query:", err)
		return false
	}

	return true
}