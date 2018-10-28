package staff

import (
	"../../fastJSON"
	"github.com/gorilla/websocket"
	"log"
	"database/sql"
	"errors"
	"encoding/json"
	"../common"
)

func auth(msg *fastjson.Value, conn *websocket.Conn) error {
	id, login, err := getUserInfo(conn, msg)
	if err != nil {
		log.Println("----[" + conn.RemoteAddr().String() + "]Error in get user info: ", err)
		return err
	}

	//поиск в уже заходивших
	oldConn, ok := getFromStaff(id)
	if ok {
		pers := common.StaffCon[oldConn]
		delete(common.StaffCon, oldConn)
		common.StaffCon[conn] = pers
		sendStaffAuthOk(conn)
		log.Println("----[" + conn.RemoteAddr().String() + "] Login as " + login)
		return nil
	}

	mytables, err := getMyTablesFromDB(id)
	if err != nil {
		log.Println("----[" + conn.RemoteAddr().String() + " " + login +
			" ]Error in get tables from database: ", err)
		common.SendError(conn, common.CommandStaffAuth, common.ErrorDBProblem)
		return err
	}

	err = addUserInOnline(conn,id,login, mytables)
	if err == nil {
		log.Println("----[" + conn.RemoteAddr().String() + "] Login as " + login)
	}

	return err
}

func getUserInfo(conn *websocket.Conn, msg *fastjson.Value) (id int, login string, err error) {
	login, pass, ok := getLoginAndPass(msg)
	if !ok {
		log.Println("----[" + conn.RemoteAddr().String() + "]Error in get login and password")
		common.SendError(conn, common.CommandStaffAuth, common.ErrorWrongCommandStructure)
		return -1,"", errors.New("Wrong command structure")
	}

	id, err = getStaffId(login, pass)
	if err != nil {
		log.Println("----[" + conn.RemoteAddr().String() + "]Error uncorrect login or password: ", err)
		if id == 0 {
			common.SendError(conn, common.CommandStaffAuth, common.ErrorWrongUser)
			return -1,"",err
		} else {
			common.SendError(conn, common.CommandStaffAuth, common.ErrorDBProblem)
			return -1,"",err
		}
	}

	return id, login, nil
}

func getLoginAndPass(msg *fastjson.Value) (login string, pass string, ok bool) {
	ok = msg.Exists("data")
	if !ok {
		return "", "", false
	}
	data := msg.Get("data")

	login = data.GetString("login")
	if login == "" {
		return "", "", false
	}

	pass = data.GetString("pass")
	if pass == "" {
		return "", "", false
	}

	return login, pass, true
}

func getStaffId(login string, pass string) (id int, err error) {
	db, err := sql.Open(common.DBConfig.UsedDatabase, common.DBConnStr)
	if err != nil {
		return -1, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT ID FROM "staff" WHERE login = $1 AND pass=$2`, login, pass)
	if err != nil {
		return -1, err
	}
	var userID sql.NullInt64
	if rows.Next() {
		rows.Scan(&userID)
		return int(userID.Int64), nil
	}
	return 0, errors.New("Can't find this user")
}

func addUserInOnline(conn *websocket.Conn, id int, login string, mytables []int) error {
	common.StaffCon[conn] = common.StaffInfo{
		Id:     id,
		Login:  login,
		Tables: mytables,
	}

	for _, tableID := range common.StaffCon[conn].Tables{
		tables := common.Tables[tableID]
		tables.Staff[id] = id
		common.Tables[tableID] = tables
	}

	return sendStaffAuthOk(conn)
}

func sendStaffAuthOk(conn *websocket.Conn) error {
	answer := common.Response{
		Command: common.CommandStaffAuth,
		Status:  true,
		Data:    nil,
	}

	jsonAnser, err := json.Marshal(answer)
	if err != nil {
		log.Println("----[" + conn.RemoteAddr().String() +"]Error in marshal Response:", err)
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonAnser)
	if err != nil {
		log.Println("----[" + conn.RemoteAddr().String() +"]Error in sending message:", err)
		return err
	}
	return nil
}

func getMyTablesFromDB(id int) (res []int, err error) {
	db, err := sql.Open(common.DBConfig.UsedDatabase, common.DBConnStr)
	if err != nil {
		return nil,err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id_table FROM staff_tables WHERE id_staff = $1`, id)
	if err != nil {
		return nil, err
	}
	res = []int{}
	id = 0
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		res = append(res, id)
	}

	return res,nil

}

func getFromStaff(id int) (conn *websocket.Conn, ok bool) {
	for conn, pers := range common.StaffCon {
		if pers.Id == id {
			return conn, true
		}
	}

	return nil, false
}
