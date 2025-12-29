// "go test" suppresses standard output unless the test fails.
// If you still want to print output use the -v option like this: "go test -v ./..."

package request

import (
	"fmt"
	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
	"io"
	"log"
	// "strings"
	"testing"
)

type chunkReader struct {
	data         string
	chunkInBytes int
	index        int
}

func (r *chunkReader) Read(p []byte) (int, error) {
	bytesToCopy := len(r.data) - r.index
	if bytesToCopy <= 0 {
		// bytesToCopy shouldn't be less than 0
		return 0, io.EOF
	}
	if bytesToCopy > r.chunkInBytes {
		bytesToCopy = r.chunkInBytes
	}
	n := copy(p, r.data[r.index:r.index+bytesToCopy])
	r.index += n
	return n, nil
}

func quickTest(r *chunkReader) {
	b := make([]byte, 1024)
	n, err := r.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("read", n, "bytes", string(b))
}

func TestRequestLineParse(t *testing.T) {
	fmt.Println("Hi!")
	// Test: Good GET Request line
	// r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	reader := chunkReader{
		data:         "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		chunkInBytes: 16,
	}
	for range 3 {
		quickTest(&reader)
	}

	// r, err := RequestFromReader(&reader)
	// require.NoError(t, err)
	// require.NotNil(t, r)
	// assert.Equal(t, "GET", r.RequestLine.Method)
	// assert.Equal(t, "/", r.RequestLine.RequestTarget)
	// assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// // Test: Good GET Request line with path
	// r, err := RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	// require.NoError(t, err)
	// require.NotNil(t, r)
	// assert.Equal(t, "GET", r.RequestLine.Method)
	// assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	// assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// // Test: Invalid number of parts in request line
	// _, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	// require.Error(t, err)
}
