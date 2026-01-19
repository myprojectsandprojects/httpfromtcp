package response

import (
	"fmt"
	"io"

	"boot.theprimeagen.tv/internal/headers"
)

type StatusCode int

const (
	StatusCode_200 StatusCode = 200
	StatusCode_400 StatusCode = 400
	StatusCode_500 StatusCode = 500
)

func (code StatusCode) String() string {
	var s string
	switch code {
	case StatusCode_200:
		s = "OK"
	case StatusCode_400:
		s = "Bad Request"
	case StatusCode_500:
		s = "Internal Server Error"
	default:
		panic("unknown status code")
	}
	return s
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := fmt.Sprintf("HTTP/1.1 %v %v\r\n", int(statusCode), statusCode)
	_, err := w.Write([]byte(statusLine)) //@ Does it write everything?
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.Create()
	h["content-length"] = fmt.Sprint(contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		header := fmt.Sprintf("%v: %v\r\n", k, v)
		_, err := w.Write([]byte(header)) //@ Does it write everything?
		if err != nil {
			return err
		}
	}
	return nil
}
