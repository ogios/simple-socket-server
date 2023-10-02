package normal

import (
	"bufio"
	"fmt"
	"net"

	"github.com/ogios/sutils"
)

type Conn struct {
	Raw    net.Conn
	Si     sutils.SBodyIN
	So     sutils.SBodyOUT
	Reader *bufio.Reader
	Type   string
}

func (c *Conn) Close() error {
	return c.Raw.Close()
}

func (c *Conn) GetType(max_length int) (string, error) {
	length, err := c.Si.Next()
	if err != nil {
		return "", err
	}
	if length > max_length {
		return "", fmt.Errorf("type too long: %d", length)
	}
	if length < 1 {
		return "", fmt.Errorf("type is empty")
	}

	b, err := c.Si.GetSec()
	if err != nil {
		return "", err
	}
	return string(b), err
}
