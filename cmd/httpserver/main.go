package main

import (
	"boot.theprimeagen.tv/internal/headers"
	"boot.theprimeagen.tv/internal/request"
	"boot.theprimeagen.tv/internal/response"
	"boot.theprimeagen.tv/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	hs := headers.Create()
	hs.Set("Connection", "close")
	hs.Set("Content-Type", "text/html")

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

		w.WriteStatusLine(response.StatusCode_400)
		w.WriteHeaders(hs)
		w.WriteBody([]byte(body))
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

		w.WriteStatusLine(response.StatusCode_500)
		w.WriteHeaders(hs)
		w.WriteBody([]byte(body))
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

		w.WriteStatusLine(response.StatusCode_200)
		w.WriteHeaders(hs)
		w.WriteBody([]byte(body))
	}
}
