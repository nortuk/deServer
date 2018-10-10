package server

import (
	"github.com/gorilla/websocket"
	"../fastJSON"
	"errors"
)

func getTableInfo(conn *websocket.Conn, msg *fastjson.Value) {

}

func getValue(msg *fastjson.Value) (val int, err error) {
	if !msg.Exists("data") {
		return 0, errors.New("Don't exist data")
	}
	data := msg.Get("data")
	if !data.Exists("value") {
		return 0, errors.New("Don't exist value")
	}
	val = data.GetInt("value")
	if err != nil {
		return 0, err
	}

	return val,err
}
