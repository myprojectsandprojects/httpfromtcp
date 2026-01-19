package request

import (
	"bytes"
	"errors"
	"io"
	"log"
	"strconv"

	"boot.theprimeagen.tv/internal/headers"
	// "strings"
)

const (
	reqStatusInitialized int = iota
	reqStatusParsingHeaders
	reqStatusParsingBody
	reqStatusDone
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte

	status        int
	contentLength int
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := Request{
		status:  reqStatusInitialized,
		Headers: headers.Create(),
	}

	buf := make([]byte, 1024)
	ri := 0
	pi := 0
	for req.status != reqStatusDone {
		if ri == len(buf) {
			log.Panicln("Request too big for our buffer.")
		}

		n, err := reader.Read(buf[ri:])
		if err != nil {
			return nil, err
		}
		ri += n

		for {
			bytesParsed, parseErr := req.parse(buf[pi:ri])
			if bytesParsed == 0 {
				if parseErr != nil {
					return nil, parseErr
				}

				break
			}

			pi += bytesParsed
		}
	}

	return &req, nil
}

// if it encounters an error, n == 0
func (r *Request) parse(d []byte) (int, error) {
	var n int
	var err error

	switch r.status {
	case reqStatusInitialized:
		var reqLine *RequestLine
		n, reqLine, err = parseRequestLine(d)
		if reqLine != nil {
			r.RequestLine = *reqLine
			r.status = reqStatusParsingHeaders
		}
	case reqStatusParsingHeaders:
		var done bool
		n, done, err = r.Headers.Parse(d)
		if done == true {
			if r.Headers.Get("Content-Length") == "" {
				r.status = reqStatusDone
			} else {
				r.contentLength, err = strconv.Atoi(r.Headers["content-length"])
				if err != nil {
					n = 0
				} else {
					r.status = reqStatusParsingBody
				}
			}
		}
	case reqStatusParsingBody:
		if len(d) < r.contentLength {
			n = 0
			err = nil
		} else if len(d) > r.contentLength {
			n = 0
			err = errors.New("too much data in body")
		} else {
			r.Body = d[:r.contentLength]
			// r.Body = make([]byte, contentLength)
			// copy(r.Body, d)

			r.status = reqStatusDone
		}
	case reqStatusDone:
		n = 0
		err = nil
	default:
		log.Panic("Unknown request status")
	}

	return n, err
}

func parseRequestLine(reqBytes []byte) (int, *RequestLine, error) {
	var method, target, httpVersion []byte

	i := bytes.Index(reqBytes, []byte{'\r', '\n'})
	if i == -1 {
		return 0, nil, nil
	}

	startLine := reqBytes[:i]

	parts := bytes.Split(startLine, []byte{' '})

	if len(parts) != 3 {
		return 0, nil, errors.New("Wrong number of elements in start-line!")
	}

	method = parts[0]
	target = parts[1]
	httpVersion = parts[2]

	// I think we should make sure that it's one of the allowed methods instead of whether it contains capital alphabetic characters only (?)
	if string(method) != "GET" && string(method) != "POST" {
		return 0, nil, errors.New("Unknown HTTP method!")
	}

	versionParts := bytes.Split(httpVersion, []byte{'/'})

	if len(versionParts) != 2 || string(versionParts[0]) != "HTTP" || string(versionParts[1]) != "1.1" {
		return 0, nil, errors.New("Wrong or malformed HTTP version!")
	}

	httpVersion = versionParts[1]

	return i + 2, &RequestLine{
		Method:        string(method),
		RequestTarget: string(target),
		HttpVersion:   string(httpVersion),
	}, nil
}
