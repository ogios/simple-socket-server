package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/ogios/sutils"
)

func Client() {
	conn, err := net.Dial("tcp", ":15001")
	if err != nil {
		panic(err)
	}
	so := sutils.NewSBodyOUT()
	so.AddBytes([]byte("test"))
	so.AddBytes([]byte("456789"))
	so.AddBytes([]byte("ayayo"))
	so.WriteTo(conn)
	si := sutils.NewSBodyIn(bufio.NewReader(conn))
	b, err := si.GetSec()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// b, err = si.GetSec()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(b))
}
