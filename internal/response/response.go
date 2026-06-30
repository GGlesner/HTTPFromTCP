// Package response responses
package response

import (
	"fmt"
	"io"

	"HTTPFromTCP/internal/headers"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func WriteStatusLine(
	w io.Writer,
	statusCode StatusCode,
) error {
	codeToMessage := map[StatusCode]string{
		OK:                  "OK",
		BadRequest:          "Bad Request",
		InternalServerError: "Internal Server Error",
	}
	message, ok := codeToMessage[statusCode]
	if !ok {
		message = ""
	}
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\n", statusCode, message)
	_, err := w.Write([]byte(statusLine))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = fmt.Sprintf("%d", contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func WriteHeaders(
	w io.Writer,
	h headers.Headers,
) error {
	for k, v := range h {
		_, err := w.Write([]byte(k + ": " + v + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
