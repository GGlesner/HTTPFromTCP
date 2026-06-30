package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"HTTPFromTCP/internal/request"
	"HTTPFromTCP/internal/response"
	"HTTPFromTCP/internal/server"
)

const port = 42069

func defaultHandler(
	w *response.Writer,
	req *request.Request,
) {
	target := req.RequestLine.RequestTarget
	var status response.StatusCode
	message := ""
	h1 := ""
	switch target {
	case "/yourproblem":
		status = response.BadRequest
		h1 = "Bad Request"
		message = "Your request honestly kinda sucked."
	case "/myproblem":
		status = response.InternalServerError
		h1 = "Internal Server Error"
		message = "Okay, you know what? This one is on me."
	default:
		status = response.OK
		h1 = "Success!"
		message = "Your request was an absolute banger."
	}
	if err := w.WriteStatusLine(status); err != nil {
		log.Println(err)
		return
	}
	h := response.GetDefaultHeaders(0)
	codeToMessage := map[response.StatusCode]string{
		response.OK:                  "OK",
		response.BadRequest:          "Bad Request",
		response.InternalServerError: "Internal Server Error",
	}
	body := fmt.Appendf(nil,
		`<html>
			<head>
				<title>%d %s</title>
			</head>
			<body>
				<h1>%s</h1>
				<p>%s</p>
			</body>
		</html>`,
		status, codeToMessage[status],
		h1,
		message,
	)
	h.Set("Content-Length", fmt.Sprintf("%d", len(body)))
	if err := w.WriteHeaders(h); err != nil {
		log.Println(err)
		return
	}
	_, err := w.WriteBody(body)
	if err != nil {
		log.Println(err)
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
