package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"boot.theprimeagen.tv/internal/headers"
	"boot.theprimeagen.tv/internal/request"
	"boot.theprimeagen.tv/internal/response"
	"boot.theprimeagen.tv/internal/server"
)

// const port = 1

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started at address:", server.Listener.Addr().String())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	h := headers.Create()
	h.Set("Connection", "close")

	prefix := "/httpbin/"
	if postfix, found := strings.CutPrefix(req.RequestLine.RequestTarget, prefix); found {
		h.Set("Transfer-Encoding", "chunked")

		err := w.WriteStatusLine(response.StatusCode_200)
		if err != nil {
			panic("failed")
		}

		url := fmt.Sprintf("https://httpbin.org/%s", postfix)
		resp, err := http.Get(url)
		if err != nil {
			panic("failed")
		}

		contentType := resp.Header["Content-Type"]
		if len(contentType) != 1 {
			panic("unexpected Content-Type header from the server")
		}
		h.Set("Content-Type", contentType[0])
		// h.Set("Content-Type", "text/plain")

		err = w.WriteHeaders(h)
		if err != nil {
			panic("failed")
		}

		buf := make([]byte, 1024)
		for {
			n, err := resp.Body.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				panic("failed to read for some unknown reason")
			}

			time.Sleep(time.Second * 3)
			_, err = w.WriteChunkedBody(buf[:n])
			if err != nil {
				panic("failed")
			}
		}
		_, err = w.WriteChunkedBodyDone()
		if err != nil {
			panic("failed")
		}

		return
	}

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		body :=
			`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>

</html>`

		h.Set("Content-Type", "text/html")
		err := w.WriteStatusLine(response.StatusCode_400)
		if err != nil {
			panic("failure")
		}
		err = w.WriteHeaders(h)
		if err != nil {
			panic("failure")
		}
		_, err = w.WriteBody([]byte(body))
		if err != nil {
			panic("failure")
		}
	case "/myproblem":
		body :=
			`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`

		h.Set("Content-Type", "text/html")
		err := w.WriteStatusLine(response.StatusCode_500)
		if err != nil {
			panic("failure")
		}
		err = w.WriteHeaders(h)
		if err != nil {
			panic("failure")
		}
		_, err = w.WriteBody([]byte(body))
		if err != nil {
			panic("failure")
		}
	default:
		body :=
			`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

		if req.RequestLine.RequestTarget == "/favicon.ico" {
			fmt.Println("started serving favicon.ico")
		}
		h.Set("Content-Type", "text/html")
		err := w.WriteStatusLine(response.StatusCode_200)
		if err != nil {
			panic("failure")
		}
		err = w.WriteHeaders(h)
		if err != nil {
			panic("failure")
		}
		_, err = w.WriteBody([]byte(body))
		if err != nil {
			panic("failure")
		}
		if req.RequestLine.RequestTarget == "/favicon.ico" {
			fmt.Println("finished serving favicon.ico")
		}
	}
}
