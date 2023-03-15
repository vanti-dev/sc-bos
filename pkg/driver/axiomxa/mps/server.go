package mps

import (
	"bufio"
	"bytes"
	"net"
	"sort"
	"sync"

	"go.uber.org/zap"
)

type OnMessageFunc = func(data []byte)

// Server receives message port packets and sends each to OnMessage for processing.
type Server struct {
	OnMessage OnMessageFunc
	Logger    *zap.Logger
}

func NewServer(handler OnMessageFunc) *Server {
	return &Server{OnMessage: handler}
}

// Serve blocks accepting connections on lis.
// lis will be closed when this returns, close lis to stop the server.
func (s *Server) Serve(lis net.Listener) error {
	if s.OnMessage == nil {
		panic("no OnMessage specified, server will do nothing")
	}

	lis = &onceCloseListener{Listener: lis}
	defer lis.Close()

	for {
		rw, err := lis.Accept()
		if err != nil {
			return err
		}

		go s.serveConn(rw)
	}
}

func (s *Server) serveConn(rw net.Conn) {
	lineReader := bufio.NewScanner(rw)
	for lineReader.Scan() {
		msg := lineReader.Bytes()
		s.logMessage(rw, msg)
		s.OnMessage(msg)
	}

	// todo: handle any errors, or errors we don't expect
	// lineReader.Err() ...
}

func (s *Server) logMessage(rw net.Conn, msg []byte) {
	if s.Logger == nil {
		return
	}
	s.Logger.Debug("message port received", zap.Stringer("remoteAddr", rw.RemoteAddr()), zap.ByteString("msg", msg))
}

// onceCloseListener wraps a net.Listener, protecting it from
// multiple Close calls.
type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() { oc.closeErr = oc.Listener.Close() }

// MapPrefix is an OnMessageFunc that forwards the data to the best matching key of routes, or calls fallback.
func MapPrefix(routes map[string]OnMessageFunc, fallback OnMessageFunc) OnMessageFunc {
	type handler struct {
		p  []byte
		do OnMessageFunc
	}
	handlers := make([]handler, 0, len(routes))
	for prefix, do := range routes {
		handlers = append(handlers, handler{p: []byte(prefix), do: do})
	}
	sort.Slice(handlers, func(i, j int) bool {
		return len(handlers[i].p) > len(handlers[j].p)
	})
	return func(data []byte) {
		for _, h := range handlers {
			if bytes.HasPrefix(data, h.p) {
				h.do(data)
				return
			}
		}
		fallback(data)
	}
}
