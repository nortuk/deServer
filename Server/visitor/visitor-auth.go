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

func auth(msg *fastjson.Value, conn *websocket.Conn) error {
	imei, tableID, err := getLoginInfo(msg)
	if err != nil {
		common.SendError(conn, common.CommandVisitorAuth, common.ErrorWrongCommandStructure)
		return err
	}

	common.TablesMutex.Lock()
	table, ok := common.Tables[tableID]
	if !ok {
		common.TablesMutex.Unlock()
		common.SendError(conn, common.CommandVisitorAuth, common.ErrorWrongCommandStructure)
		return errors.New("tableID uncorrect")
	}
	table.Visitors = append(table.Visitors, imei)
	common.Tables[tableID] = table
	common.TablesMutex.Unlock()

	common.VisitorsConnMutex.Lock()
	common.VisitorsConn[conn] = common.VisitorInfo{
		IMEI:    imei,
		TableID: tableID,
	}
	common.VisitorsConnMutex.Unlock()

	err = sendVisitorAuthOk(conn)
	if err != nil {
		common.SendError(conn, common.CommandVisitorAuth, common.ErrorWrongCommandStructure)
		return err
	}

	return nil
}

func getLoginInfo(msg *fastjson.Value) (imei string, tableID int, err error) {
	if !msg.Exists("IMEI") {
		return "", 0, errors.New("IMEI doesn't exists")
	}
	imei = msg.GetString("IMEI")
	if checkIMEI(imei) {
		return "", 0, errors.New("IMEI doesn't correct")
	}
	if !msg.Exists("tableID") {
		return "", 0, errors.New("tableID doesn't exists")
	}
	tableID = msg.GetInt("tableID")

	return imei, tableID, nil
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
		log.Println("----["+conn.RemoteAddr().String()+"]Error in marshal Response:", err)
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonAnser)
	if err != nil {
		log.Println("----["+conn.RemoteAddr().String()+"]Error in sending message:", err)
		return err
	}
	return nil
}
