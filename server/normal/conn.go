package normal

import (
	"bufio"
	"net"
)

type Conn struct {
	Raw    net.Conn
	Reader *bufio.Reader
	Type   string
}

func (c *Conn) Close() error {
	return c.Raw.Close()
}
