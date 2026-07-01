// Package response responses
package response

import (
	"errors"
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

type WriterState string

const (
	RequestLine WriterState = "request-line"
	Headers     WriterState = "headers"
	Body        WriterState = "body"
)

type Writer struct {
	IOWriter    io.Writer
	writerState WriterState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		IOWriter:    w,
		writerState: RequestLine,
	}
}

func (w *Writer) WriteStatusLine(
	statusCode StatusCode,
) error {
	if w.writerState != RequestLine {
		return errors.New("trying to write request-line in the wrong state: " + string(w.writerState))
	}
	codeToMessage := map[StatusCode]string{
		OK:                  "OK",
		BadRequest:          "Bad Request",
		InternalServerError: "Internal Server Error",
	}
	message, ok := codeToMessage[statusCode]
	if !ok {
		message = ""
	}
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, message)
	_, err := w.IOWriter.Write([]byte(statusLine))
	if err != nil {
		return err
	}
	w.writerState = Headers
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = fmt.Sprintf("%d", contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/html"
	return h
}

func (w *Writer) WriteHeaders(
	h headers.Headers,
) error {
	if w.writerState != Headers {
		return errors.New("trying to write headers in the wrong state: " + string(w.writerState))
	}
	for k, v := range h {
		_, err := w.IOWriter.Write([]byte(k + ": " + v + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.IOWriter.Write([]byte("\r\n"))
	w.writerState = Body
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != Body {
		return 0, errors.New("trying to write body in the wrong state: " + string(w.writerState))
	}
	return w.IOWriter.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	nBytes := 0
	n, err := w.WriteBody(fmt.Appendf(nil, "%X\r\n", len(p)))
	if err != nil {
		return 0, err
	}
	nBytes += n
	n, err = w.WriteBody(p)
	if err != nil {
		return 0, err
	}
	nBytes += n
	n, err = w.WriteBody([]byte("\r\n"))
	if err != nil {
		return 0, err
	}
	nBytes += n
	return nBytes, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := w.WriteBody([]byte("0\r\n\r\n"))
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	w.writerState = Headers
	defer func() { w.writerState = Body }()
	return w.WriteHeaders(h)
}
