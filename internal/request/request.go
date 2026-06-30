// Package request Parses requests
package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"HTTPFromTCP/internal/headers"
)

const (
	crlf       = "\r\n"
	BufferSize = 8
)

type ParsingState int

const (
	ParsingRequestLine ParsingState = iota
	ParsingHeaders
	ParsingBody
	Done
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte

	state          ParsingState
	bodyLengthRead int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{
		state:   ParsingRequestLine,
		Headers: headers.NewHeaders(),
		Body:    make([]byte, 0),
	}
	buffer := make([]byte, BufferSize)
	bytesRead := 0
	for req.state < Done {
		if bytesRead >= len(buffer) {
			newBuf := make([]byte, len(buffer)*2)
			copy(newBuf, buffer)
			buffer = newBuf
		}
		n, err := reader.Read(buffer[bytesRead:])
		if err == io.EOF {
			if req.state < Done {
				return nil, fmt.Errorf("incomplete request, in state %d, read bytes on EOF: %d", req.state, bytesRead)
			}
			break
		}
		if err != nil {
			return nil, err
		}
		bytesRead += n
		bytesParsed, err := req.parse(buffer[:bytesRead])
		if err != nil {
			return nil, err
		}
		if bytesParsed > 0 {
			copy(buffer, buffer[bytesParsed:bytesRead])
			bytesRead -= bytesParsed
		}
	}
	return req, nil
}

func parseRequestLine(buffer []byte) (*RequestLine, int, error) {
	i := bytes.Index(buffer, []byte(crlf))
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
	if len(httpAndVersion) != 2 || httpAndVersion[0] != "HTTP" || httpAndVersion[1] != "1.1" {
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
	bytesParsed := 0
	for r.state < Done {
		n, err := r.parseSingle(buffer[bytesParsed:])
		if err != nil {
			return 0, err
		}
		bytesParsed += n
		if n == 0 {
			break
		}
	}
	return bytesParsed, nil
}

func (r *Request) parseSingle(buffer []byte) (int, error) {
	switch r.state {
	case ParsingRequestLine:
		rql, n, err := parseRequestLine(buffer)
		if err != nil {
			return 0, err
		} else if n == 0 {
			return n, nil
		}
		r.RequestLine = *rql
		r.state = ParsingHeaders
		return n, nil
	case ParsingHeaders:
		n, done, err := r.Headers.Parse(buffer)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = ParsingBody
		}
		return n, nil
	case ParsingBody:
		length, ok := r.Headers.Get("Content-Length")
		if !ok {
			r.state = Done
			return len(buffer), nil
		}
		numBytes, err := strconv.Atoi(length)
		if err != nil {
			return 0, errors.New("invalid Content-Length value: " + length)
		}
		r.Body = append(r.Body, buffer...)
		r.bodyLengthRead += len(buffer)
		if numBytes < r.bodyLengthRead {
			return 0, fmt.Errorf("incorrect Content-Length value: expected %s, actual %d", length, len(r.Body))
		} else if r.bodyLengthRead == numBytes {
			r.state = Done
		}
		return len(buffer), err
	case Done:
		return 0, fmt.Errorf("error: trying to read from a done state")
	default:
		return 0, errors.New("unknown state")
	}
}
