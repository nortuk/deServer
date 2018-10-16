package server

import (
	"../fastJSON"
	"database/sql"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"./common"
)

type personType int

const (
	isVisitorType personType = 0
	isStaffType   personType = 1
	unknownType   personType = 2
)

func getAuthType(msg *fastjson.Value) (res personType) {
	ok := isVisitor(msg)
	if ok {
		return isVisitorType
	}

	ok = isStaff(msg)
	if ok {
		return isStaffType
	}

	return unknownType
}

func isStaff(jsonMsg *fastjson.Value) bool {
	ok := jsonMsg.Exists("command")
	if !ok {
		return false
	}

	if jsonMsg.GetString("command") == common.CommandStaffAuth {
		return true
	}

	return false
}

func isVisitor(jsonMsg *fastjson.Value) bool {
	ok := jsonMsg.Exists("type")
	if !ok {
		return false
	}

	if jsonMsg.GetString("type") == common.CommandVisitorAuth {
		return true
	}

	return false
}

func staffAuth(msg *fastjson.Value, conn *websocket.Conn) bool {
	login, pass, ok := getLoginAndPass(msg)
	if !ok {
		log.Println("Error in get login and password")
		common.SendError(conn, common.CommandStaffAuth, common.ErrorWrongCommandStructure)
		return false
	}

	id, err := getStaffId(login, pass)
	if err != nil {
		log.Println("Error uncorrect login or password: ", err)
		if id == 0 {
			common.SendError(conn, common.CommandStaffAuth, common.ErrorWrongUser)
			return false
		} else {
			common.SendError(conn, common.CommandStaffAuth, common.ErrorDBProblem)
			return false
		}
	}

	mytables, err := getMyTablesFromDB(id)
	if err != nil {
		common.SendError(conn, common.CommandStaffAuth, common.ErrorDBProblem)
		return false
	}

	common.StaffCon[conn] = common.StaffInfo{
		Id:     id,
		Login:  login,
		Pass:   pass,
		Tables: mytables,
	}

	_, ok = common.StaffCon[conn]
	if !ok {
		log.Println("Error in adding staff")
		common.SendError(conn, common.CommandStaffAuth, common.ErrorCannotAddStaff)
		return false
	}

	sendStaffAuthOk(conn)
	return true
}

func getLoginAndPass(msg *fastjson.Value) (login string, pass string, ok bool) {
	ok = msg.Exists("data")
	if !ok {
		log.Println("Error in staff parse: data don't exists")
		return "", "", false
	}
	data := msg.Get("data")

	login = data.GetString("login")
	if login == "" {
		log.Println("Error in staff parse: login don't exists")
		return "", "", false
	}

	pass = data.GetString("pass")
	if pass == "" {
		log.Println("Error in staff parse: password uncorrect")
		return "", "", false
	}

	return login, pass, true
}

func getStaffId(login string, pass string) (id int, err error) {
	db, err := sql.Open(common.DBConfig.UsedDatabase, common.DBConnStr)
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		return -1, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT ID FROM "staff" WHERE login = $1 AND pass=$2`, login, pass)
	if err != nil {
		log.Println("Error in request execution:", err)
		return -1, err
	}
	var userID sql.NullInt64
	if rows.Next() {
		rows.Scan(&userID)
		return int(userID.Int64), nil
	}
	return 0, errors.New("Can't find this user")
}

func sendStaffAuthOk(conn *websocket.Conn) {
	answer := common.Response{
		Command: common.CommandStaffAuth,
		Status:  true,
		Data:    nil,
	}

	jsonAnser, err := json.Marshal(answer)
	if err != nil {
		log.Println("Error in marshal Response:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonAnser)
	if err != nil {
		log.Println("Error in sending message:", err)
		return
	}
}

func getMyTablesFromDB(id int) (res []int, err error) {
	db, err := sql.Open(common.DBConfig.UsedDatabase, common.DBConnStr)
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		return nil,err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id_table FROM staff_tables WHERE id_staff = $1`, id)
	if err != nil {
		log.Println("Error in request execution:", err)
		return nil, err
	}
	res = []int{}
	id = 0
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Println("Error: in reading my tables: ", err)
			return nil, err
		}
		res = append(res, id)
	}

	return res,nil

}