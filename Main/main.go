package main

import (
	_ "svn.cloudserver.ru/Rendall_server/fastJSON"
	"time"
	_ "../Handlers"
)

func main() {
	for {
		time.Sleep(1000 * time.Second)
	}
	return
}
