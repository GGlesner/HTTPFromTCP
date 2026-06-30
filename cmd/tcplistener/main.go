package main

import (
	"fmt"
	"log"
	"net"

	"HTTPFromTCP/internal/request"
)

const (
	port       = ":42069"
	bufferSize = 8
)

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error listening for TCP traffic on port %s: %s\n", port, err.Error())
	}
	defer listener.Close()
	fmt.Printf("Listening for TCP traffic on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error listening to port %s: %s\n", port, err.Error())
		}
		fmt.Println("Accepted connection from ", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err.Error())
		}

		rql := req.RequestLine
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", rql.Method, rql.RequestTarget, rql.HttpVersion)
		fmt.Println("Connection to ", conn.RemoteAddr(), " closed")

		headers := req.Headers
		fmt.Print("Headers:\n")
		for key, value := range headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

	}
}
