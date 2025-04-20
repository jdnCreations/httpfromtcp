package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/jdnCreations/httpfromtcp/internal/request"
	"github.com/jdnCreations/httpfromtcp/internal/response"
)
type Handler func(w *response.Writer, req *request.Request) 

type Server struct {
	listener net.Listener
  handler Handler
	closed atomic.Bool
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
	
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	w := response.NewWriter(conn)
  req, err := request.RequestFromReader(conn) 
	if err != nil {
		w.WriteStatusLine(response.StatusBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}

	s.handler(w, req)
}

func (s *Server) listen() {

	for !s.closed.Load() {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}

	server := &Server{
		listener: listener,
    handler: handler,
	}

	go func() {
		server.listen()
	}()

	return server, nil 
}