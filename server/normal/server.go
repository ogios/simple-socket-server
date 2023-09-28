package normal

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/ogios/simple-socket-server/config"
	"github.com/ogios/simple-socket-server/log"
)

type Server struct {
	Listener      net.Listener
	Addr          string
	typeCallbacks map[string]TypeCallback
	cond          sync.Cond
}

type TypeCallback func(*Conn) error

func NewSocketServer() (*Server, error) {
	l, err := net.Listen("tcp", config.GLOBAL_CONFIG.Server.Addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		Listener:      l,
		Addr:          config.GLOBAL_CONFIG.Server.Addr,
		cond:          *sync.NewCond(&sync.Mutex{}),
		typeCallbacks: map[string]TypeCallback{},
	}, nil
}

// infinent serv loop
func (s *Server) Serv() error {
	log.Info(nil, "Server serving...")
	for {
		conn, err := s.Listener.Accept()
		log.Info(nil, "connected: %s", conn.RemoteAddr().String())
		if err != nil {
			return err
		}
		go s.Process(conn)
	}
}

// for loop runs callbacks
func (s *Server) Process(conn net.Conn) {
	// close Connection
	defer func() {
		log.Info(nil, "Connection <%s> closed", conn.RemoteAddr().String())
		conn.Close()
	}()

	// catch error
	defer func() {
		if err := recover(); err != nil {
			log.Error(nil, "Unexpected Connection process error: %s", err)
		}
	}()

	// read type until \n
	log.Debug(nil, "Connection <%s> start processing", conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)
	t, err := reader.ReadString(0xa)
	if err != nil {
		panic(err)
	}
	if len(t) < 1 {
		panic("Type is empty")
	}
	t = t[:len(t)-1]

	// get callback and execute
	s.cond.L.Lock()
	process, ok := s.typeCallbacks[t]
	s.cond.L.Unlock()
	if !ok {
		panic(fmt.Sprintf("Unknow type: %s", t))
	}
	err = process(&Conn{
		Raw:    conn,
		Reader: reader,
		Type:   t,
	})
	if err != nil {
		log.Error(nil, "Process error: %s", err)
	}
	log.Debug(nil, "Connection <%s> process done", conn.RemoteAddr().String())
}

// add callbacks for certain type
func (s *Server) AddTypeCallback(t string, callback TypeCallback) (okk bool) {
	if strings.Contains(t, "\n") {
		log.Error(nil, "type does not support string with \\n")
		return false
	}
	s.cond.L.Lock()
	s.typeCallbacks[t] = callback
	s.cond.L.Unlock()
	return true
}
