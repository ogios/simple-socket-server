package main

import (
	"fmt"
	"os"

	"golang.org/x/exp/slog"

	_ "github.com/ogios/simple-socket-server/config"
	"github.com/ogios/simple-socket-server/log"
	"github.com/ogios/simple-socket-server/server/normal"
)

func test() {
	log.SetLevel(slog.LevelDebug)
	server, err := normal.NewSocketServer()
	fmt.Println("server created")
	if err != nil {
		panic(err)
	}
	server.AddTypeCallback("push", func(conn *normal.Conn) error {
		// read and print every single byte till the end
		fmt.Printf("Type: %s\n", conn.Type)
		b, readerr := conn.Si.GetSec()
		conn.Close()
		if readerr != nil {
			return readerr
		}
		fmt.Printf("%d ", b)
		return nil
	})
	if err := server.Serv(); err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) >= 1 {
		test()
	}
}
