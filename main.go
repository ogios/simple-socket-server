package main

import (
	"fmt"

	_ "transfer-go/config"
	_ "transfer-go/log"
	"transfer-go/server/normal"
)

func main() {
	server, err := normal.NewSocketServer()
	fmt.Println("server created")
	if err != nil {
		panic(err)
	}
	if err := server.Serv(); err != nil {
		panic(err)
	}
}
