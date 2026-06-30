// Package server serves
package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"HTTPFromTCP/internal/request"
	"HTTPFromTCP/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(
	port int,
	handler Handler,
) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: listener,
		handler:  handler,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(
	conn net.Conn,
) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		he := &HandlerError{
			StatusCode: response.BadRequest,
			Message:    err.Error(),
		}
		_, err := he.Write(conn)
		if err != nil {
			log.Println(err)
		}
	}
	buf := bytes.NewBuffer(make([]byte, 0))
	he := s.handler(buf, req)
	if he != nil {
		_, err := he.Write(conn)
		if err != nil {
			log.Println(err)
		}
		return
	}
	err = response.WriteStatusLine(conn, response.OK)
	if err != nil {
		log.Println(err)
	}
	header := response.GetDefaultHeaders(buf.Len())
	err = response.WriteHeaders(conn, header)
	if err != nil {
		log.Println(err)
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Println(err)
	}
}

func (he *HandlerError) Write(
	w io.Writer,
) (int, error) {
	err := response.WriteStatusLine(w, he.StatusCode)
	if err != nil {
		return 0, err
	}
	message := []byte(he.Message)
	err = response.WriteHeaders(w, response.GetDefaultHeaders(len(message)))
	if err != nil {
		return 0, err
	}
	return w.Write(message)
}
