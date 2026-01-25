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

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		header := fmt.Sprintf("%v: %v\r\n", k, v)
		_, err := w.Write([]byte(header)) //@ We can't assume that it writes all (source: Perplexity)
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

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) Writer {
	return Writer{
		writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	_, err := fmt.Fprintf(w.writer, "HTTP/1.1 %v %v\r\n", int(statusCode), statusCode)
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for k, v := range headers {
		_, err := fmt.Fprintf(w.writer, "%v: %v\r\n", k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := fmt.Fprintf(w.writer, "Content-Length: %v\r\n\r\n%s", len(p), p)
	if err != nil {
		return 0, err
	}

	return n, nil
}
