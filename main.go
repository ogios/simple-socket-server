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
		for {
			b, readerr := conn.Reader.ReadByte()
			if readerr != nil {
				if readerr.Error() == "EOF" {
					fmt.Print("\n")
					conn.Close()
					return nil
				} else {
					return readerr
				}
			}
			fmt.Printf("%d ", b)
		}
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
