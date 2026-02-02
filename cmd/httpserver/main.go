package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	// "time"

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
		h.Set("Trailer", "X-Content-Length, X-Content-SHA256")

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

		var content []byte
		buf := make([]byte, 64)
		for {
			n, err := resp.Body.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				panic("failed to read for some unknown reason")
			}

			// time.Sleep(time.Second * 3)
			_, err = w.WriteChunkedBody(buf[:n])
			if err != nil {
				panic("failed") //@ Obviously the server shouldn't panic (crash) when the client disconnects midway
			}

			content = append(content, buf[:n]...)
		}
		_, err = w.WriteChunkedBodyDone()
		if err != nil {
			panic("failed")
		}

		contentLen := strconv.Itoa(len(content))
		checksum := fmt.Sprintf("%x", sha256.Sum256(content))
		err = w.WriteTrailers(h, contentLen, checksum)
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
	case "/video":
		f := "/home/eero/go-test/httpfromtcp/assets/vim.mp4"
		body, err := os.ReadFile(f)
		if err != nil {
			panic("failed to read the file")
		}
		fmt.Printf("file size: %v\n", len(body))

		h.Set("Content-Type", "video/mp4")
		// h.Set("Content-Disposition", "attachment; filename=\"vim.mp4\"") // Tells browsers that the response should be treated as a file download, not rendered inline (by built-in players or viewers)

		err = w.WriteStatusLine(response.StatusCode_200)
		if err != nil {
			panic("failure")
		}
		err = w.WriteHeaders(h)
		if err != nil {
			panic("failure")
		}
		_, err = w.WriteBody(body)
		if err != nil {
			log.Print(err)
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
	}
}
