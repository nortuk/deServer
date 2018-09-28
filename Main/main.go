package main

import (
	"../Server"
	"log"
)

func main() {
	err := server.Start("server.json","database.json")
	if err != nil {
		log.Print("Error in server initialization:", err)
	}

	return
}
