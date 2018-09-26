package main

import (
	_ "svn.cloudserver.ru/Rendall_server/fastJSON"
	"../Config"
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"sync"
	"flag"
	"math/rand"
	"time"
	"database/sql"
)

type (
	connection struct {
		id              int
		protocolVersion int
		sessionId       string
		ws              *websocket.Conn
		sync.RWMutex
	}

	wsMsg struct {
		msgType int
		message string
	}
)

var (
	addr *string
	chars = []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0123456789")
	socketCounter = 0
	db *sql.DB
	charsLen = len(chars)
	seededRand *rand.Rand
	sessionIdLen = 30
	wsConnections struct {
					  sync.RWMutex
					  list map[int]*connection
	}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func initRand() {
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randGenString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rune(chars[rand.Intn(charsLen)])
	}
	return string(b)
}

func messageReader(conn *connection) {
	locked := false
	for {
		if locked {conn.Unlock()}
		locked = false
		msgType, msg, err := conn.ws.ReadMessage()
		if err != nil {
			log.Println("ReadMessage ERROR:", err)
			break
		}
		log.Println("new msg:", string(msg), string(msgType), " from: ", conn.id)
		conn.Lock()
		locked = true
		switch msgType {
		case websocket.TextMessage:
			log.Println("Get text msg")
			// parse message
		case websocket.CloseMessage:
			log.Println("Closing connection")
			// delete conn
		default:
			log.Println("Unexpected message, not supported!")
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("UPGRADE ERROR", err)
		return
	}
	log.Println("Connection:", conn.RemoteAddr())
	user := connection{
		id: socketCounter,
		ws: conn,
		sessionId: randGenString(sessionIdLen),
	}
	socketCounter++
	messageReader(&user)
}

func main() {
	wsConnections.list = make(map[int]*connection)
	log.Println("start init ws server")
	initRand()

	log.Println("init db connection")
	db, err := Config.ReadDatabaseConfig("../Config/database.json")
	if err != nil {
		log.Println("error in init db connection:", err, db)
	}
	log.Println("Success db connection, db:", db)

	var server Config.Server
	server, err = Config.ReadServerConfig("../Config/server.json")
	if err != nil {
		log.Println("error in parse server config:", err, server)
	}
	value := server.Ip + ":" + server.Port
	addr = flag.String("addr", value, "localhost")

 	http.HandleFunc("/", handler)
	log.Println("start server on addr:", *addr)
	log.Fatal(http.ListenAndServe(*addr,nil))
}
