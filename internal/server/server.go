package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/jdnCreations/httpfromtcp/internal/request"
	"github.com/jdnCreations/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
  handler Handler
	closed atomic.Bool
}

type HandlerError struct {
  StatusCode response.StatusCode
  Message string
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
	
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (he *HandlerError) Write(w io.Writer) {
  response.WriteStatusLine(w, he.StatusCode)
	messageBytes := []byte(he.Message)
  headers := response.GetDefaultHeaders(len(messageBytes))
  response.WriteHeaders(w, headers)
	w.Write(messageBytes)
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

  req, err := request.RequestFromReader(conn) 
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message: err.Error(),
		}
		hErr.Write(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	hErr := s.handler(buf, req)
	if hErr != nil {
		hErr.Write(conn)
		return
	}

	b := buf.Bytes() 
	response.WriteStatusLine(conn, response.StatusOK) 
  headers := response.GetDefaultHeaders(len(b))
  response.WriteHeaders(conn, headers)
	conn.Write(b)
	return	
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