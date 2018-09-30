package main

import (
	"../Server"
	"log"
)

func main() {
	err := server.Start("server.json","database.json")
	if err != nil {
		log.Println("Error in server initialization:", err)
	}

	return
}
