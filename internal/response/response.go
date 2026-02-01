package response

import (
	"fmt"
	"io"
	"strconv"
	"strings"

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

type Writer struct {
	writer            io.Writer
	wroteEndOfHeaders bool
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

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	totalWritten := 0
	if !w.wroteEndOfHeaders {
		n, err := fmt.Fprintf(w.writer, "\r\n")
		if err != nil {
			return 0, err
		}
		w.wroteEndOfHeaders = true
		totalWritten += n
	}

	hexStr := strconv.FormatInt(int64(len(p)), 16)
	n, err := fmt.Fprintf(w.writer, "%s\r\n%s\r\n", hexStr, p)
	if err != nil {
		return 0, err
	}
	totalWritten += n

	return totalWritten, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := fmt.Fprintf(w.writer, "0\r\n\r\n")
	if err != nil {
		return 0, err
	}

	return n, nil
}

// map[string]string // where key is the trailer name and the value is the trailer value. Feels cleaner
func (w *Writer) WriteTrailers(h headers.Headers, values ...string) error {
	names := strings.Split(h.Get("Trailer"), ", ")
	if len(names) != len(values) {
		panic("fail")
	}

	// Perplexity argues that this should not be sent (which makes sense to me):
	// _, err := fmt.Fprintf(w.writer, "0\r\n")
	// if err != nil {
	// 	panic("fail")
	// }

	for i, name := range names {
		_, err := fmt.Fprintf(w.writer, "%v: %v\r\n", name, values[i])
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w.writer, "\r\n")
	if err != nil {
		panic("fail")
	}

	return nil
}
