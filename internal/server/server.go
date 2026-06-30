// Package server serves
package server

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: listener,
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
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	response := strings.Join([]string{
		"HTTP/1.1 200 OK",
		"Content-Type: text/plain",
		"Content-Length: 13",
		"",
		"Hello World!\n",
	}, "\r\n")
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Fatalf("error writing response: %v", err)
	}
}
