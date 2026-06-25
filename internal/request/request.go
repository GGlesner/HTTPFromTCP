// Package request Parses requests
package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	var req *Request
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return req, err
	}
	requestLine, err := parseRequestLine(bytes)
	if err != nil {
		return req, err
	}
	req = &Request{
		RequestLine: requestLine,
	}
	return req, nil
}

func parseRequestLine(bytes []byte) (RequestLine, error) {
	rql := RequestLine{}
	line := string(bytes)
	i := 0
	for ; i < len(line)-1; i++ {
		if line[i:i+2] == "\r\n" {
			break
		}
	}
	parts := strings.Split(line[:i], " ")
	if len(parts) != 3 {
		return rql, fmt.Errorf("expected 3 parts in request line, got: %d", len(parts))
	}
	for _, chr := range parts[0] {
		if !unicode.IsLetter(chr) || !unicode.IsUpper(chr) {
			return rql, fmt.Errorf("unknow http method: %s", parts[0])
		}
	}
	httpAndVersion := strings.Split(parts[2], "/")
	if len(httpAndVersion) < 2 || httpAndVersion[0] != "HTTP" || httpAndVersion[1] != "1.1" {
		return rql, fmt.Errorf("expected last part to be HTTP/1.1, got: %s", parts[2])
	}
	rql = RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpAndVersion[1],
	}
	return rql, nil
}
