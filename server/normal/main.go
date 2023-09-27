package normal

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"

	"transfer-go/config"
	"transfer-go/log"
)

type Server struct {
	Listener      net.Listener
	Addr          string
	typeCallbacks map[string][]func(net.Conn) error
	cond          sync.Cond
}

func NewSocketServer() (*Server, error) {
	l, err := net.Listen("tcp", config.SysConfig.Server.Addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		Listener: l,
		Addr:     config.SysConfig.Server.Addr,
		cond:     *sync.NewCond(&sync.Mutex{}),
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
	defer conn.Close()
	defer func() {
		if err := recover(); err != nil {
			log.Error(nil, "Connection process error: %s", err)
		}
	}()
	log.Debug(nil, "Connection <%s> start processing", conn.RemoteAddr().String())
	conn.Close()
	reader := bufio.NewReader(conn)
	t, err := reader.ReadString(0xa)
	if err != nil {
		panic(err)
	}
	s.cond.L.Lock()
	processes, ok := s.typeCallbacks[t]
	s.cond.L.Unlock()
	if !ok {
		panic(fmt.Sprintf("Unknow type: %s", t))
	}
	for _, process := range processes {
		err := process(conn)
		if err != nil {
			log.Error(nil, "Process error: %s", err)
		}
	}
}

// add callbacks for certain type
func (s *Server) AddTypeCallback(t string, callback func(net.Conn) error) (okk bool) {
	if strings.Contains(t, "\n") {
		log.Error(nil, "type does not support string with \\n")
		return false
	}
	s.cond.L.Lock()
	_, ok := s.typeCallbacks[t]
	if !ok {
		s.typeCallbacks[t] = []func(net.Conn) error{callback}
	} else {
		s.typeCallbacks[t] = append(s.typeCallbacks[t], callback)
	}
	s.cond.L.Unlock()
	return true
}
