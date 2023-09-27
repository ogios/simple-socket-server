package test

import (
	"bufio"
	"fmt"
	"math"
	"net"

	"transfer-go/api/request/types"
	"transfer-go/log"

	"golang.org/x/exp/slog"
)

func GetLen(reader *bufio.Reader) int {
	readlen, err := reader.ReadBytes(0)
	rawlen := readlen[:len(readlen)-1]
	if err != nil {
		return 0
	}
	total := 0
	for index, b := range rawlen {
		pow := float64(len(rawlen) - 1 - index)
		feat := int(math.Pow(255, pow))
		total += int(b) * feat
	}
	return total
}

func Process(conn net.Conn) {
	defer conn.Close()
	defer func() {
		if err := recover(); err != nil {
			log.Error(nil, "Connection process error: %s", err)
		}
	}()
	slog.Debug("Connection <%s> start processing", conn.RemoteAddr().String())
	conn.Close()
	reader := bufio.NewReader(conn)
	// t, err := reader.ReadString(0xa)
	// if err != nil {
	// 	panic(err)
	// }
	t := GetLen(reader)
	fmt.Println("head: ", t)
	switch t {
	case types.Push:

	case types.Fetch:

	default:
		panic(fmt.Sprintf("Unknow type: %d", t))
	}
}
