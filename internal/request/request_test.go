// "go test" suppresses standard output unless the test fails.
// If you still want to print output use the -v option like this: "go test -v ./..."

package request

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	// "strings"
	"testing"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos > len(cr.data) {
		log.Panic("Something went completely haywire!")
	}
	if cr.pos == len(cr.data) {
		return 0, io.EOF
	}

	// endIndex := cr.pos + cr.numBytesPerRead
	// if endIndex > len(cr.data) {
	// 	endIndex = len(cr.data)
	// }
	endIndex := min((cr.pos + cr.numBytesPerRead), len(cr.data))
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	// fmt.Printf("copied %v bytes: %q\n", n, p[:n])
	return n, nil
}

func TestRequestLineParse(t *testing.T) {
	// reader := chunkReader{
	// 	data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	// }
	// reader.numBytesPerRead = len(reader.data)

	// var wholeThing []byte
	// b := make([]byte, 1024)
	// for {
	// 	n, err := reader.Read(b)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		t.Fatal("error occured")
	// 	}
	// 	fmt.Printf("%v bytes: %q\n", n, b[:n])
	// 	wholeThing = append(wholeThing, b[:n]...)
	// }
	// fmt.Printf("whole thing: %q\n", wholeThing)

	// Test: Good GET Request line
	reader := chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(&reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	reader = chunkReader{
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}
	r, err = RequestFromReader(&reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line
	reader = chunkReader{
		data:            "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 16,
	}
	r, err = RequestFromReader(&reader)
	require.Nil(t, r)
	require.Error(t, err)

	// Test: wrong method
	reader = chunkReader{
		data: "get /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	}
	reader.numBytesPerRead = len(reader.data)
	r, err = RequestFromReader(&reader)
	require.Nil(t, r)
	require.Error(t, err)
}

func TestRequestParse(t *testing.T) {
	// Test: Standard Headers
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "localhost:42069", r.Headers["host"])
	assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
	assert.Equal(t, "*/*", r.Headers["accept"])

	// Test: Malformed Header
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.Error(t, err)

	// Test: Empty headers
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	assert.Equal(t, 0, len(r.Headers))

	// Test: Duplicate Headers
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\nSuperhero: Batman\r\nSome:other header\r\nsuperhero: Superman\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	assert.Equal(t, "Batman, Superman", r.Headers["superhero"])

	// // Test: Duplicate Headers (only '\n')
	// reader = &chunkReader{
	// 	data:            "GET / HTTP/1.1\nSuperhero: Batman\nSome:other header\nsuperhero: Superman\n\n",
	// 	numBytesPerRead: 3,
	// }
	// r, err = RequestFromReader(reader)
	// require.NoError(t, err)
	// require.NotNil(t, r)
	// assert.Equal(t, "GET", r.RequestLine.Method)
	// assert.Equal(t, "/", r.RequestLine.RequestTarget)
	// assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	// assert.Equal(t, "Batman, Superman", r.Headers["superhero"])

	// // Test: Missing end of headers
	// reader = &chunkReader{
	// 	data:            "GET / HTTP/1.1\r\nSuperhero: Batman",
	// 	numBytesPerRead: 3,
	// }
	// r, err = RequestFromReader(reader)
	// require.NoError(t, err)
	// require.NotNil(t, r)
	// assert.Equal(t, "GET", r.RequestLine.Method)
	// assert.Equal(t, "/", r.RequestLine.RequestTarget)
	// assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestBodyParse(t *testing.T) {
	// Test: Standard Body
	reader := &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 13\r\n" +
			"\r\n" +
			"hello world!\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "hello world!\n", string(r.Body))

	// Test: Body shorter than reported content length
	reader = &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 20\r\n" +
			"\r\n" +
			"partial content",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.Nil(t, r)
	require.Error(t, err)

	//@ Add some more test cases. Here are the names of some of my additional test cases:
	//
	// "Standard Body" (valid)
	// "Empty Body, 0 reported content length" (valid)
	// "Empty Body, no reported content length" (valid)
	// "Body shorter than reported content length" (should error)
	// "No Content-Length but Body Exists" (shouldn't error; we're assuming Content-Length will be present if a body exists)
}
