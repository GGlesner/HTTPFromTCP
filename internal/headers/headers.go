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

const crlf = "\r\n"

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
	key := string(data[:colonIdx])
	if key != strings.Trim(key, " ") {
		return 0, false, errors.New("no leading nor trailing spaces")
	}
	value := strings.Trim(string(data[colonIdx+1:i]), " ")
	if len(value) == 0 {
		return 0, false, errors.New("no value")
	}
	h[key] = value
	return i + 2, false, nil
}
