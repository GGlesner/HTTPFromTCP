// Package server serves
package server

import (
	"fmt"
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

type Handler func(w *response.Writer, req *request.Request)

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
	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Println(err)
		return
	}
	s.handler(w, req)
	// err = w.WriteStatusLine(response.OK)
	// if err != nil {
	// 	log.Println(err)
	// }
	// header := response.GetDefaultHeaders(buf.Len())
	// err = w.WriteHeaders(header)
	// if err != nil {
	// 	log.Println(err)
	// }
	// _, err = conn.Write(buf.Bytes())
	// if err != nil {
	// 	log.Println(err)
	// }
}
