package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := min(cr.pos+cr.numBytesPerRead, len(cr.data))
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}

func TestRequestLineParse(t *testing.T) {
	// Test: Good GET Request line
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1024,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, r.RequestLine.Method, "GET")
	assert.Equal(t, r.RequestLine.RequestTarget, "/")
	assert.Equal(t, r.RequestLine.HttpVersion, "1.1")

	// Test: Good GET Request line with path
	reader.data = "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	reader.pos = 0
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, r.RequestLine.Method, "GET")
	assert.Equal(t, r.RequestLine.RequestTarget, "/coffee")
	assert.Equal(t, r.RequestLine.HttpVersion, "1.1")

	// Test: Good POST request line with path
	reader.data = "POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	reader.pos = 0
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, r.RequestLine.Method, "POST")
	assert.Equal(t, r.RequestLine.RequestTarget, "/coffee")
	assert.Equal(t, r.RequestLine.HttpVersion, "1.1")

	// Test: Invalid number of parts in request line
	reader.data = "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	reader.pos = 0
	_, err = RequestFromReader(reader)
	require.Error(t, err)

	// Test: Invalid method request line
	reader.data = "get /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	reader.pos = 0
	_, err = RequestFromReader(reader)
	require.Error(t, err)

	// Test: Invalid version in request line
	reader.data = "GET /coffee HTTP/1.2\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	reader.pos = 0
	_, err = RequestFromReader(reader)
	require.Error(t, err)
}
