package normal

import (
	"bufio"
	"net"

	"github.com/ogios/sutils"
)

type Conn struct {
	Raw    net.Conn
	Si     *sutils.SBodyIN
	So     *sutils.SBodyOUT
	Reader *bufio.Reader
	Type   string
}

func (c *Conn) Close() error {
	return c.Raw.Close()
}
