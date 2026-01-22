package response

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

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

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		header := fmt.Sprintf("%v: %v\r\n", k, v)
		_, err := w.Write([]byte(header)) //@ Does it write everything?
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		panic("boom")
	}

	return nil
}

type Response struct {
	Bytes      bytes.Buffer
	StatusCode StatusCode
	Headers    headers.Headers
}

func New() *Response {
	return &Response{
		Headers: headers.Create(),
	}
}

func (resp *Response) SetStatusCode(code StatusCode) {
	resp.StatusCode = code
}

func (resp *Response) SetDefaultHeaders() {
	resp.Headers.Set("connection", "close")
}

func (resp *Response) Body(body string) {
	// statusLine := fmt.Sprintf("HTTP/1.1 %v %v\r\n", int(resp.StatusCode), resp.StatusCode)
	// _, err := resp.Bytes.Write([]byte(statusLine)) //@ Does it write everything?
	// if err != nil {
	// 	panic("boom")
	// }
	_, err := fmt.Fprintf(&resp.Bytes, "HTTP/1.1 %v %v\r\n", int(resp.StatusCode), resp.StatusCode)
	if err != nil {
		panic("boom")
	}

	// contentLen := fmt.Sprintf("%v", len(body))
	contentLen := strconv.Itoa(len(body))
	resp.Headers.Set("Content-Length", contentLen)

	err = WriteHeaders(&resp.Bytes, resp.Headers)
	if err != nil {
		panic("boom")
	}

	_, err = resp.Bytes.Write([]byte(body))
	if err != nil {
		panic("boom")
	}
}
