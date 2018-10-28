package main

import (
	"../Server"
	"log"
	"os"
	"time"
	"strconv"
)

func main() {
	t := time.Now()
	name := "log-"
	name += strconv.Itoa(t.Day()) + "." + t.Month().String() +
		"." + strconv.Itoa(t.Year()) + ".log"
	f, err := os.OpenFile(name, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	err = server.Start("server.json","database.json")
	if err != nil {
		log.Println("Error in server initialization:", err)
	}

	return
}
