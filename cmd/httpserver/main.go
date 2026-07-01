package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"HTTPFromTCP/internal/headers"
	"HTTPFromTCP/internal/request"
	"HTTPFromTCP/internal/response"
	"HTTPFromTCP/internal/server"
)

const (
	port        = 42069
	DefaultHTML = `<html>
			<head>
				<title>%d %s</title>
			</head>
			<body>
				<h1>%s</h1>
				<p>%s</p>
			</body>
		</html>`
)

func handler(
	w *response.Writer,
	req *request.Request,
) {
	target := req.RequestLine.RequestTarget
	if strings.HasPrefix(target, "/httpbin/") {
		httpbinHandler(w, req)
		return
	}
	switch target {
	case "/yourproblem":
		handler400(w, req)
	case "/myproblem":
		handler500(w, req)
	case "/video":
		videoHandler(w, req)
	default:
		handler200(w, req)
	}
}

func handler200(
	w *response.Writer,
	_ *request.Request,
) {
	if err := w.WriteStatusLine(response.OK); err != nil {
		log.Println(err)
		return
	}
	h := response.GetDefaultHeaders(0)
	body := fmt.Appendf(nil,
		DefaultHTML,
		response.OK, "OK",
		"Success!",
		"Your request was an absolute banger.",
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

func handler400(
	w *response.Writer,
	_ *request.Request,
) {
	if err := w.WriteStatusLine(response.BadRequest); err != nil {
		log.Println(err)
		return
	}
	h := response.GetDefaultHeaders(0)
	body := fmt.Appendf(nil,
		DefaultHTML,
		response.BadRequest, "Bad Request",
		"Bad Request",
		"Your request honestly kinda sucked.",
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

func handler500(
	w *response.Writer,
	_ *request.Request,
) {
	if err := w.WriteStatusLine(response.InternalServerError); err != nil {
		log.Println(err)
		return
	}
	h := response.GetDefaultHeaders(0)
	body := fmt.Appendf(nil,
		DefaultHTML,
		response.InternalServerError, "Internal Server Error",
		"Internal Server Error",
		"Okay, you know what? This one is on me.",
	)
	h.Set("Content-Length", fmt.Sprintf("%d", len(body)))
	if err := w.WriteHeaders(h); err != nil {
		log.Println(err)
		return
	}
	_, err := w.WriteBody(body)
	if err != nil {
		log.Println(err)
		return
	}
}

func videoHandler(
	w *response.Writer,
	req *request.Request,
) {
	video, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		log.Println(err)
		handler500(w, req)
		return
	}
	if err := w.WriteStatusLine(response.OK); err != nil {
		log.Println(err)
		return
	}
	h := headers.NewHeaders()
	h["Content-Type"] = "video/mp4"
	h["Content-Length"] = fmt.Sprintf("%d", len(video))
	if err = w.WriteHeaders(h); err != nil {
		log.Println(err)
		return
	}
	if _, err = w.WriteBody(video); err != nil {
		log.Println(err)
		return
	}
}

func httpbinHandler(
	w *response.Writer,
	req *request.Request,
) {
	req.RequestLine.RequestTarget = strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
	delete(req.Headers, "Content-Length")
	req.Headers["Transfer-Encoding"] = "chunked"
	res, err := http.Get("https://httpbin.org" + req.RequestLine.RequestTarget)
	if err != nil {
		log.Println(err)
		handler400(w, req)
		return
	}
	if err := w.WriteStatusLine(response.OK); err != nil {
		log.Println(err)
		return
	}

	h := response.GetDefaultHeaders(0)
	delete(h, "Content-Length")
	h["Transfer-Encoding"] = "chunked"
	h["Trailers"] = "X-Content-Sha256, X-Content-Length"
	err = w.WriteHeaders(h)
	if err != nil {
		log.Println(err)
		return
	}

	p := make([]byte, 1024)
	body := make([]byte, 0)
	for {
		n, err := res.Body.Read(p)
		if err == io.EOF {
			_, err = w.WriteBody([]byte("0\r\n"))
			if err != nil {
				log.Println(err)
				return
			}
			break
		} else if err != nil {
			log.Println(err)
			return
		}
		_, err = w.WriteChunkedBody(p[:n])
		if err != nil {
			log.Println(err)
			return
		}
		body = append(body, p[:n]...)
	}

	hash := sha256.Sum256(body)
	t := headers.NewHeaders()
	t["X-Content-Sha256"] = hex.EncodeToString(hash[:])
	t["X-Content-Length"] = fmt.Sprintf("%d", len(body))
	err = w.WriteTrailers(t)
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	server, err := server.Serve(port, handler)
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
