// Package headers: ...
package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

const (
	crlf = "\r\n"
)

func (h Headers) Get(key string) (string, bool) {
	if val, ok := h[strings.ToLower(key)]; ok {
		return val, true
	} else {
		return "", false
	}
}

func (h Headers) Set(key string, val string) {
	h[key] = val
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	i := bytes.Index(data, []byte(crlf))
	if i < 0 {
		return 0, false, nil
	}
	if i == 0 {
		return 2, true, nil
	}
	colonIdx := bytes.Index(data[:i], []byte(":"))
	if colonIdx < 0 {
		return 0, false, errors.New("no ':'")
	}
	if colonIdx == 0 {
		return 0, false, errors.New("no field-name")
	}
	tchar := "!#$%&'*+-.^_`|~"
	key := strings.ToLower(string(data[:colonIdx]))
	if key != strings.TrimSpace(key) {
		return 0, false, errors.New("no leading nor trailing spaces")
	}
	if strings.ContainsAny(key, "\t") {
		return 0, false, errors.New("no whitespaces")
	}
	for _, r := range key {
		chr := string(r)
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && !strings.Contains(tchar, chr) {
			return 0, false, errors.New("invalid character: " + chr)
		}
	}
	value := strings.TrimSpace(string(data[colonIdx+1 : i]))
	if len(value) == 0 {
		return 0, false, errors.New("no value")
	}
	if previous, ok := h[key]; !ok {
		h[key] = value
	} else {
		h[key] = previous + ", " + value
	}
	return i + 2, false, nil
}
