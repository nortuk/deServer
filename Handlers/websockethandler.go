package handler

import (
	//"svn.cloudserver.ru/Rendall_server/fastJSON"
	"database/sql"
	"math/rand"
	"../Config"
	"sync"
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"time"
	"flag"
)

type (
	Connection struct {
		id              int
		sessionId       string
		ws              *websocket.Conn
		sync.RWMutex
	}

	WsMsg struct {
		msgType int
		message string
	}
)


var (
	chars = []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0123456789")
	socketCounter = 0
	seededRand *rand.Rand
	db *sql.DB
	addr *string
	charsLen = len(chars)
	sessionIdLen = 30
	wsConnections struct {
			 sync.RWMutex
			 list map[int]*Connection
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

func init () {
	var server Config.Server
	server, err := Config.ReadServerConfig("../Config/server.json")
	if err != nil {
		log.Println("error in parse server config:", err, server)
	}
	value := server.Ip + ":" + server.Port
	addr = flag.String("addr", value, "localhost")
	wsConnections.list = make(map[int]*Connection)
	log.Println("start init ws server")
	initRand()

	log.Println("init db connection")
	db, err := Config.ReadDatabaseConfig("../Config/database.json")
	if err != nil {
		log.Println("error in init db connection:", err, db)
	}
	log.Println("Success db connection, db:", db)


	http.HandleFunc("/", WsHandler)
	log.Println("start server on addr:", *addr)
	log.Fatal(http.ListenAndServe(*addr,nil))
}

func messageReader(conn *Connection) {
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
			// ниже пример работы с fastJSON парсером.
			// инициализация
			/*
			p := fastjson.Parser{}
			// получение *fastJSON.Value
			value, err := p.Parse(string(msg))
			if err != nil {
				log.Println("error in parse fastJSON", err)
				continue
			}
			log.Println("FASTJSON VALUE:", value)
			// извлечение данных из *fastJSON.Value по ключу (ключ-всегда строка)
			// есть GetInt and so on...
			str := value.GetStringBytes("key")
			log.Println(string(str))
			// TODO parse messages from client
			*/
		case websocket.CloseMessage:
			log.Println("Closing connection")
			// delete conn
		default:
			log.Println("Unexpected message, not supported!")
		}
	}
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("UPGRADE ERROR", err)
		return
	}
	log.Println("Connection:", conn.RemoteAddr())
	user := Connection{
		id: socketCounter,
		ws: conn,
		sessionId: randGenString(sessionIdLen),
	}
	socketCounter++
	messageReader(&user)
}

