package main

import (
	"bufio"
	"fmt"
	"net"
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
	server.AddTypeCallback("push", func(conn net.Conn, reader *bufio.Reader) error {
		// read and print every single byte till the end
		for {
			b, readerr := reader.ReadByte()
			if readerr != nil {
				if readerr.Error() == "EOF" {
					conn.Close()
					return nil
				} else {
					return readerr
				}
			}
			fmt.Println(b)
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
