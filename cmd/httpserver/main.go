package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"HTTPFromTCP/internal/request"
	"HTTPFromTCP/internal/response"
	"HTTPFromTCP/internal/server"
)

const port = 42069

func defaultHandler(w io.Writer, req *request.Request) *server.HandlerError {
	target := req.RequestLine.RequestTarget
	var status response.StatusCode
	message := ""
	switch target {
	case "/yourproblem":
		status = response.BadRequest
		message = "Your problem is not my problem\n"
	case "/myproblem":
		status = response.InternalServerError
		message = "Woopsie, my bad\n"
	default:
		_, err := w.Write([]byte("All good, frfr\n"))
		if err != nil {
			log.Println(err)
		}
		return nil
	}
	return &server.HandlerError{
		StatusCode: status,
		Message:    message,
	}
}

func main() {
	server, err := server.Serve(port, defaultHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port ", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
