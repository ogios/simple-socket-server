package main

import (
	"os"

	"github.com/ogios/simple-socket-server/server/normal"
	"golang.org/x/exp/slog"
)

var PASS string = "456789"

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// Level: slog.LevelDebug,
	})))
	server, err := normal.NewSocketServer(":15001")
	if err != nil {
		panic(err)
	}

	Fill(server)

	go Client()

	if err = server.Serv(); err != nil {
		panic(err)
	}
}
