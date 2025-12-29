package request

import (
	"io"
	// "fmt"
	"bytes"
	"errors"
	// "log"
	// "strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	reqBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var reqLine *RequestLine
	reqLine, err = parseRequestLine(reqBytes)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *reqLine,
	}, nil
}

func parseRequestLine(reqBytes []byte) (*RequestLine, error) {
	var method, target, httpVersion []byte

	i := bytes.Index(reqBytes, []byte{'\r', '\n'})
	if i == -1 {
		return nil, errors.New("Malformed request!")
	}

	startLine := reqBytes[:i]

	parts := bytes.Split(startLine, []byte{' '})

	if len(parts) != 3 {
		return nil, errors.New("Wrong number of elements in start-line!")
	}

	method = parts[0]
	target = parts[1]
	httpVersion = parts[2]

	versionParts := bytes.Split(httpVersion, []byte{'/'})

	if len(versionParts) != 2 || string(versionParts[0]) != "HTTP" || string(versionParts[1]) != "1.1" {
		return nil, errors.New("Wrong or malformed HTTP version!")
	}

	httpVersion = versionParts[1]

	return &RequestLine{
		Method:        string(method),
		RequestTarget: string(target),
		HttpVersion:   string(httpVersion),
	}, nil
}
