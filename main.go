package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		for line := range parsLines(conn) {
			fmt.Println(line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), " closed")
	}
}

func parsLines(file io.ReadCloser) <-chan string {
	lines := make(chan string)
	buffer := make([]byte, bufferSize)
	line := ""
	sendLines := func() {
		defer file.Close()
		defer close(lines)
		for {
			n, err := file.Read(buffer)
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatal(err)
			}
			parts := strings.Split(string(buffer[:n]), "\n")
			m := len(parts)
			parts[0] = line + parts[0]
			for i := range m - 1 {
				lines <- parts[i]
			}
			line = parts[m-1]
		}
	}
	go sendLines()
	return lines
}
