package normal

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/ogios/sutils"

	"github.com/ogios/simple-socket-server/log"
)

type Server struct {
	Listener      net.Listener
	Addr          string
	typeCallbacks map[string]TypeCallback
	sms           []StartMiddleware
	ems           []EndMiddleware
	cond          sync.Cond
	MaxTypeLength int
}

type TypeCallback func(*Conn) error
type StartMiddleware func(*Conn) error
type EndMiddleware func(conn *Conn, err any)

func NewSocketServer(addr string) (*Server, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		Listener:      l,
		Addr:          l.Addr().String(),
		cond:          *sync.NewCond(&sync.Mutex{}),
		typeCallbacks: map[string]TypeCallback{},
		sms:           make([]StartMiddleware, 0),
		ems:           make([]EndMiddleware, 0),
		MaxTypeLength: 1024,
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
		go s.process(conn)
	}
}

// add callbacks for certain type
func (s *Server) AddTypeCallback(t string, callback TypeCallback) (okk bool) {
	if strings.Contains(t, "\n") {
		log.Error(nil, "type gradle clean test --testsdoes not support string with \\n")
		return false
	}
	s.cond.L.Lock()
	s.typeCallbacks[t] = callback
	s.cond.L.Unlock()
	return true
}

// add middleware at the start of a connection
func (s *Server) AddMiddlewareOnStart(call StartMiddleware) {
	s.sms = append(s.sms, call)
}

// add middleware at the end of process
func (s *Server) AddMiddlewareOnEnd(call EndMiddleware) {
	s.ems = append(s.ems, call)
}

func (s *Server) execute(conn *Conn) {
	// catch error
	defer func() {
		err := recover()
		if err != nil {
			log.Error(nil, "Process error: %s", err)
			defer conn.Close()
		}
		for _, em := range s.ems {
			em(conn, err)
		}
	}()

	// get type
	t, err := getType(conn.Si, s.MaxTypeLength)
	if err != nil {
		panic(err)
	}
	conn.Type = t

	// execute middleware before callback
	for index, sm := range s.sms {
		err = sm(conn)
		if err != nil {
			log.Error(nil, "execute start middleware error: index-%d msg-%s", index, err)
			panic(err)
		}
	}

	// get callback and execute
	s.cond.L.Lock()
	process, ok := s.typeCallbacks[conn.Type]
	s.cond.L.Unlock()
	if !ok {
		panic(fmt.Errorf("unknow type: %s", conn.Type))
	}
	err = process(conn)
	if err != nil {
		log.Error(nil, "Process error: %s", err)
		panic(err)
	}
}

// for loop runs callbacks
func (s *Server) process(conn net.Conn) {
	defer log.Debug(nil, "Connection <%s> process done", conn.RemoteAddr().String())

	// second insurance
	defer func() {
		if err := recover(); err != nil {
			conn.Close()
			log.Error(nil, "Unexpected end middleware error: %s", err)
		}
	}()

	// create conn
	log.Debug(nil, "Connection <%s> start processing", conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)
	cc := new(Conn)
	cc.Raw = conn
	cc.Si = sutils.NewSBodyIn(reader)
	cc.So = sutils.NewSBodyOUT()
	cc.Reader = reader

	// start process
	s.execute(cc)
}

func getType(si *sutils.SBodyIN, max_length int) (string, error) {
	length, err := si.Next()
	if err != nil {
		return "", err
	}
	if length > max_length {
		return "", fmt.Errorf("type too long: %d", length)
	}
	if length < 1 {
		return "", fmt.Errorf("type is empty")
	}

	b, err := si.GetSec()
	if err != nil {
		return "", err
	}
	return string(b), err
}
