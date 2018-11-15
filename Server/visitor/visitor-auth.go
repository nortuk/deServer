package visitor

import (
	"../../fastJSON"
	"github.com/gorilla/websocket"
	"errors"
	"../common"
	"unicode"
	"encoding/json"
	"log"
)

func auth(msg *fastjson.Value, conn *websocket.Conn) error{
	if !msg.Exists("IMEI") {
		common.SendError(conn, common.CommandVisitorAuth, common.ErrorWrongCommandStructure)
		return errors.New("IMEI doesn't exists")
	}
	imei := msg.GetString("IMEI")
	if checkIMEI(imei) {
		common.SendError(conn, common.CommandVisitorAuth, common.ErrorWrongCommandStructure)
		return errors.New("IMEI doesn't correct")
	}
	if !msg.Exists("tableID") {
		common.SendError(conn, common.CommandVisitorAuth, common.ErrorWrongCommandStructure)
		return errors.New("tableID doesn't exists")
	}
	tableID := msg.GetInt("tableID")

	table, ok := common.Tables[tableID]
	if !ok {
		common.SendError(conn, common.CommandVisitorAuth, common.ErrorWrongCommandStructure)
		return errors.New("tableID uncorrect")
	}
	table.Visitors = append(table.Visitors, imei)

	common.Tables[tableID] = table

	common.VisitorsConn[conn] = common.VisitorInfo{
		IMEI: imei,
		TableID: tableID,
	}

	err := sendVisitorAuthOk(conn)
	if err != nil {
		common.SendError(conn, common.CommandVisitorAuth, common.ErrorWrongCommandStructure)
		return err
	}

	return nil
}

func checkIMEI(imei string) bool {
	if len(imei) != 18 {
		return false
	}

	for _, smb := range imei {
		if !unicode.IsDigit(rune(smb)) {
			return false
		}
	}

	return true
}

func sendVisitorAuthOk(conn *websocket.Conn) error {
	answer := common.Response{
		Command: common.CommandVisitorAuth,
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