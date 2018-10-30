package common

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type errCode int

const (
	ErrorUnknownMsg            errCode = 100
	ErrorUnknownCommandType    errCode = 101
	ErrorWrongCommandStructure errCode = 102
	ErrorWrongUser             errCode = 103
	ErrorDBProblem             errCode = 104
	ErrorProductDontExists     errCode = 105
	ErrorTableDoesnotExists    errCode = 106
	ErrorConnectionrefused 	   errCode = 107
)

func SendError(conn *websocket.Conn, command string, errorCode errCode) {
	answer := Response{
		Command: command,
		Status:  false,
		Data: DataStruct{
			"value": errorCode,
		},
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
