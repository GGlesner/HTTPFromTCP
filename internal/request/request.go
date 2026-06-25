// Package request Parses requests
package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	CRLF       = "\r\n"
	BufferSize = 1024
)

type Request struct {
	RequestLine RequestLine
	state       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{}
	buffer := make([]byte, BufferSize)
	bytesRead := 0
	bytesParsed := 0
	for req.state < 1 {
		if bytesRead >= len(buffer) {
			newBuf := make([]byte, len(buffer)*2)
			copy(newBuf, buffer)
			buffer = newBuf
		}
		n, err := reader.Read(buffer[bytesRead:])
		if err == io.EOF {
			req.state = 1
			break
		}
		if err != nil {
			return nil, err
		}
		bytesRead += n
		bytesParsed, err = req.parse(buffer[:bytesRead])
		if err != nil {
			return nil, err
		}
		if bytesParsed > 0 {
			newBuf := make([]byte, bytesRead-bytesParsed)
			copy(newBuf, buffer[bytesParsed:bytesRead])
			buffer = newBuf
			bytesRead -= bytesParsed
		}
	}
	return req, nil
}

func parseRequestLine(buffer []byte) (*RequestLine, int, error) {
	i := bytes.Index(buffer, []byte(CRLF))
	if i < 0 {
		return nil, 0, nil
	}
	line := string(buffer[:i])
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, 0, fmt.Errorf("expected 3 parts in request line, got: %d", len(parts))
	}
	for _, chr := range parts[0] {
		if !unicode.IsLetter(chr) || !unicode.IsUpper(chr) {
			return nil, 0, fmt.Errorf("unknow http method: %s", parts[0])
		}
	}
	httpAndVersion := strings.Split(parts[2], "/")
	if len(httpAndVersion) < 2 || httpAndVersion[0] != "HTTP" || httpAndVersion[1] != "1.1" {
		return nil, 0, fmt.Errorf("expected last part to be HTTP/1.1, got: %s", parts[2])
	}
	return &RequestLine{
			Method:        parts[0],
			RequestTarget: parts[1],
			HttpVersion:   httpAndVersion[1],
		},
		i + 2,
		nil
}

func (r *Request) parse(buffer []byte) (int, error) {
	switch r.state {
	case 0:
		rql, n, err := parseRequestLine(buffer)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *rql
		r.state = 1
		return n, nil
	case 1:
		return 0, fmt.Errorf("error: trying to read from a done state")
	default:
		return 0, nil
	}
}
