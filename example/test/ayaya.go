package main

import (
	"fmt"
	"strings"

	"github.com/ogios/simple-socket-server/server/normal"
)

func Fill(server *normal.Server) {
	// callback
	server.AddTypeCallback("test", func(conn *normal.Conn) error {
		err := conn.So.AddBytes([]byte("ayaya!"))
		if err != nil {
			return err
		}
		err = conn.So.WriteTo(conn.Raw)
		if err != nil {
			return err
		}
		conn.Close()
		return nil
	})

	// on connected
	server.AddMiddlewareOnStart(func(conn *normal.Conn) error {
		length, err := conn.Si.Next()
		if err != nil {
			return err
		}
		if length > 10 {
			return fmt.Errorf("passwd too long: %d", length)
		}
		p, err := conn.Si.GetSec()
		if err != nil {
			return err
		}
		if strings.Compare(string(p), PASS) != 0 {
			return fmt.Errorf("wrong passwd")
		}
		return nil
	})

	// on end or error
	server.AddMiddlewareOnEnd(func(conn *normal.Conn, e any) {
		if e != nil {
			err := []byte(fmt.Sprintf("%s", e))
			conn.So.AddBytes([]byte("error"))
			conn.So.AddBytes(err)
			conn.So.WriteTo(conn.Raw)
		} else {
			fmt.Println("\n!!!!!!!!!!!!!!!!!done!!!!!!!!!!!!!!\n")
		}
	})
}
