package staff

import (
	"../../fastJSON"
	"github.com/gorilla/websocket"
	"errors"
	"log"
	"../common"
	"database/sql"
	"encoding/json"
)

type personalInfo struct {
	FirstName string	`json:"first_name"`
	SecondName string	`json:"second_name"`
	Patronymic string	`json:"patronymic"`
	Position string		`json:"position"`
}

func gettableinfo(msg *fastjson.Value, conn *websocket.Conn)  {
	tableID, err := getTableIDFromMsg(msg)
	if err != nil{
		log.Println("[" + conn.RemoteAddr().String() +"]Wrong command structure:", err)
		common.SendError(conn,common.CommandGetTableInfo, common.ErrorWrongCommandStructure)
		return
	}


	err = sendTableInfo(conn, tableID)
	if err != nil{
		log.Println("[" + conn.RemoteAddr().String() +"]Error in sending ok command(gettableinfo):", err)
	}
}

func getTableIDFromMsg(msg *fastjson.Value) (id int, err error) {
	if !msg.Exists("data"){
		return -1, errors.New("Don't exists data")
	}
	data := msg.Get("data")

	if !data.Exists("value"){
		return -1, errors.New("Don't exists value")
	}
	id = data.GetInt("value")
	return id, nil
}

func sendTableInfo(conn *websocket.Conn, tableID int ) error {
	staff, err := getTableStaffInfo(tableID)
	if err != nil{
		return  err
	}

	table, ok := common.Tables[tableID]
	if !ok {
		return errors.New("Table doesn't exists")
	}

	answer := common.Response{
		Command: common.CommandGetTableInfo,
		Status: true,
		Data: common.DataStruct{
			"staff": staff,
			"order": table.Order,
		},
	}

	jsonAnswer, err := json.Marshal(answer)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage,jsonAnswer)
	if err != nil {
		return err
	}

	return nil
}

func getTableStaffInfo(tableID int) (staff []personalInfo, err error)  {
	db, err := sql.Open(common.DBConfig.UsedDatabase, common.DBConnStr)
	if err != nil {
		return staff, err
	}
	defer db.Close()

	table, ok := common.Tables[tableID]
	if !ok{
		return staff, errors.New("Table doesn't exists")
	}

	var firstName, secondName, partonymic, position  sql.NullString
	for _, id := range table.Staff{
		rows, err := db.Query(`SELECT inf.name, inf.surname, inf.patronymic, pos.name FROM staff_info AS inf INNER JOIN positions AS pos ON inf.id_position = pos.id WHERE inf.id = $1`, id)
		if err != nil {
			return staff,err
		}
		rows.Next()
		rows.Scan(&firstName, &secondName, &partonymic, &position)
		if err != nil {
			return staff,err
		}

		staff = append(staff, personalInfo{
			FirstName: firstName.String,
			SecondName: secondName.String,
			Patronymic: partonymic.String,
			Position: position.String,
		})
	}

	return staff, nil
}