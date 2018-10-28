package server

import (
	"database/sql"
	_"github.com/lib/pq"
	"log"
	"fmt"
	"./common"
	"strconv"
)

func initializeDBWork(path string) (err error) {
	common.DBConfig, err = loadDBCfg(path)
	if err != nil {
		log.Println("EROR: load database config:", err)
		return err
	}

	if common.DBConfig.Password == "" {
		common.DBConnStr = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
			common.DBConfig.Host, common.DBConfig.Port, common.DBConfig.User, common.DBConfig.NameDatabase)
	} else{
		common.DBConnStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			common.DBConfig.Host, common.DBConfig.Port, common.DBConfig.User,
			common.DBConfig.Password, common.DBConfig.NameDatabase)
	}

	err = checkDBConnection()
	if err != nil {
		log.Println("Error in initialize database work:", err)
		return err
	}

	return nil
}

func checkDBConnection() (err error) {
	db, err := sql.Open(common.DBConfig.UsedDatabase, common.DBConnStr)
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("Error in ping database:", err)
		return err
	}

	return nil
}

func loadTables() error {
	db, err := sql.Open(common.DBConfig.UsedDatabase, common.DBConnStr)
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		return err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, name FROM tables`)
	if err != nil {
		log.Println("Error in request execution:", err)
		return err
	}
	var id sql.NullInt64
	var name sql.NullString
	for rows.Next() {
		err = rows.Scan(&id,&name)
		if err != nil {
			log.Println("Error: in reading tables: ", err)
			return err
		}
		common.Tables[int(id.Int64)] = common.TableInfo{
			Name: name.String,
			Visitors: []string{},
			Staff: map[int]int{},
		}
	}

	return nil
}

func loadMenu() error {
	db, err := sql.Open(common.DBConfig.UsedDatabase, common.DBConnStr)
	if err != nil {
		log.Println("Error in the open connection with database:", err)
		return err
	}
	defer db.Close()

	err = getCategories(db)
	if err != nil {
		log.Println("Error in get categories:", err)
		return err
	}

	return nil
}

func getCategories(db *sql.DB) error {
	rowsCategory, err := db.Query(`SELECT id, name FROM categories`)
	if err != nil {
		log.Println("Error in request execution:", err)
		return err
	}

	var idCategory sql.NullInt64
	var nameCategory sql.NullString
	for rowsCategory.Next() {
		err = rowsCategory.Scan(&idCategory,&nameCategory)
		if err != nil {
			log.Println("Error: in reading categories: ", err)
			return err
		}

		goods, err := getGoods(db , int(idCategory.Int64))
		if err != nil {
			log.Println("Error: in reading goods: ", err)
			return err
		}

		common.Menu = append(common.Menu, common.MenuCategory{
			Name: nameCategory.String,
			Goods: goods,
		})
		log.Println("goods of:", strconv.Itoa(int(idCategory.Int64)), goods)

	}

	return nil
}

func getGoods(db *sql.DB, categoryID int) (goods []common.MenuGood, err error) {
	rowsGoods, err := db.Query(`SELECT id,name,price from goods where id_category = $1`, categoryID)
	if err != nil {
		log.Println("Error in request execution:", err)
		return nil, err
	}

	var idGoods sql.NullInt64
	var name sql.NullString
	var price sql.NullInt64
	for rowsGoods.Next() {
		err = rowsGoods.Scan(&idGoods,&name,&price)
		if err != nil {
			log.Println("Error: in reading goods: ", err)
			return nil, err
		}
		goods = append(goods, common.MenuGood{
			Id: int(idGoods.Int64),
			Name: name.String,
			Price: int(price.Int64),
		})
	}
	return goods, nil
}