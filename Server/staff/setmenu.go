package staff

import (
	"../../fastJSON"
	"github.com/gorilla/websocket"
	"../common"
	"log"
	"errors"
	"encoding/json"
)

func setmenu(msg *fastjson.Value, conn *websocket.Conn)  {
	tableID, order, err := getOrderFromMsg(msg)
	if err != nil {
		log.Println("[" + conn.RemoteAddr().String() +"]Wrong command structure:", err)
		common.SendError(conn,common.CommandSetmenu, common.ErrorWrongCommandStructure)
		return
	}

	for _, ord := range order {
		if !existsProductInMenu(ord.ProductID){
			log.Println("[" + conn.RemoteAddr().String() +"]Product don't exists in menu:", err)
			common.SendError(conn,common.CommandSetmenu, common.ErrorProductDontExists)
			return
		}
	}

	table, ok := common.Tables[tableID]
	if !ok {
		log.Println("[" + conn.RemoteAddr().String() +"]Table doesn't exists:", err)
		common.SendError(conn,common.CommandSetmenu, common.ErrorTableDoesnotExists)
		return
	}
	table.Order = order
	common.Tables[tableID] = table
	if !sendSetMenuOK(conn) {
		log.Println("[" + conn.RemoteAddr().String() +"]Error in sending ok command(settables)")
	}
}

func sendSetMenuOK(conn *websocket.Conn) bool {
	answer := common.Response{
		Command: common.CommandSetmenu,
		Status: true,
		Data: nil,
	}

	jsonAnser, err := json.Marshal(answer)
	if err != nil {
		log.Println("----[" + conn.RemoteAddr().String() +"]ERROR in marshal response:", err)
		return false
	}

	err = conn.WriteMessage(websocket.TextMessage,jsonAnser)
	if err != nil {
		log.Println("----[" + conn.RemoteAddr().String() +"]ERROR in sending message:", err)
		return false
	}

	return true
}

func getOrderFromMsg(msg *fastjson.Value) (tableID int, order []common.SubOrder, err error){
	if !msg.Exists("data") {
		return -1,nil, errors.New("Don't exist data")
	}
	data := msg.Get("data")

	if !data.Exists("tableid") {
		return -1,nil, errors.New("Don't exist tableid")
	}
	tableID= data.GetInt("tableid")

	if !data.Exists("order"){
		return -1,nil, errors.New("Don't exist goods")
	}
	ord, err := data.Get("order").Array()
	if err != nil {
		return -1,nil, err
	}
	for _, subOrd := range ord {
		if !subOrd.Exists("id"){
			return -1,nil, errors.New("Don't exist id")
		}
		id := subOrd.GetInt("id")

		if !subOrd.Exists("count"){
			return -1,nil,errors.New("Don't exist count")
		}
		count := subOrd.GetInt("count")

		order = append(order, common.SubOrder{
			ProductID: id,
			Count: count,
		})
	}
	return tableID, order, err
}

func existsProductInMenu(id int)  bool {
	for _,category := range common.Menu{
		for _,product := range category.Goods{
			if product.Id == id {
				return true
			}
		}	
	}
	return false
}