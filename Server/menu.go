package server

import (
	"github.com/gorilla/websocket"
	"database/sql"
	"../Database"
	"log"
	"strconv"
	"encoding/json"
)

type Category struct {
	Name 	string 		`json:"nameofcategory"`
	Goods 	[]Goods		`json:"goods"`
}

type Goods struct {
	Id 		int 	`json:"id"`
	Name	string 	`json:"name"`
	Price 	int		`json:"price"`
}

func getMenu(conn *websocket.Conn) {
	db, err := sql.Open(database.Cfg.UsedDatabase, database.ConnStr)
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		sendError(conn,"getmenu",err.Error())
	}
	defer db.Close()

	rowsCategory, err := db.Query(`SELECT id, name FROM categories`)
	if err != nil {
		log.Println("Error in request execution:", err)
		sendError(conn,"getmenu",err.Error())
	}
	var menu []Category
	var idCategory sql.NullInt64
	var nameCategory sql.NullString
	for rowsCategory.Next() {
		err = rowsCategory.Scan(&idCategory,&nameCategory)
		if err != nil {
			log.Println("Error: in reading tables: ", err)
			sendError(conn,"getmenu",err.Error())
		}
		goods, err := getGoods(int(idCategory.Int64))
		if err != nil {
			log.Println("Error: in reading tables: ", err)
			sendError(conn,"getmenu",err.Error())
		}
		menu = append(menu, Category{
			Name: nameCategory.String,
			Goods: goods,
		})
		log.Println("goods of:", strconv.Itoa(int(idCategory.Int64)), goods)

	}
	log.Println(menu)
	if !sendMenu(conn, menu) {
		log.Println("Error in sendnenu")
		sendError(conn,"sendmenu","Error in sendmenu")
	}
}

func getGoods(id int) ([]Goods, error){
	db, err := sql.Open(database.Cfg.UsedDatabase, database.ConnStr)
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		return nil, err
	}
	defer db.Close()

	rowsGoods, err := db.Query(`SELECT id,name,price from goods where id_category = $1`, id)
	if err != nil {
		log.Println("Error in request execution:", err)
		return nil, err
	}
	var idGoods sql.NullInt64
	var name sql.NullString
	var price sql.NullInt64
	var goods []Goods
	for rowsGoods.Next() {
		err = rowsGoods.Scan(&idGoods,&name,&price)
		if err != nil {
			log.Println("Error: in reading tables: ", err)
			return nil, err
		}
		goods = append(goods, Goods{
			Id: int(idGoods.Int64),
			Name: name.String,
			Price: int(price.Int64),
		})
	}
	return goods, nil
}

func sendMenu(conn *websocket.Conn, menu []Category) bool {
	answer := response{
		Command: "getmenu",
		Status: true,
		Data: dataStruct{
			"value": menu,
		},
	}

	jsonAnswer, err := json.Marshal(answer)
	if err != nil {
		log.Println("ERROR in marshal response:", err)
		return false
	}

	err = conn.WriteMessage(websocket.TextMessage,jsonAnswer)
	if err != nil {
		log.Println("ERROR in sending message:", err)
		return false
	}

	return true
}