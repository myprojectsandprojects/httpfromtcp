package main

import (
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
	log.Println("Server started at address:", server.Addr().String())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(res *response.Response, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		res.SetStatusCode(response.StatusCode_400)
		res.SetDefaultHeaders()
		res.Headers.Set("Content-Type", "text/plain")
		res.Body("Your problem is not my problem\n")
	case "/myproblem":
		res.SetStatusCode(response.StatusCode_500)
		res.SetDefaultHeaders()
		res.Headers.Set("Content-Type", "text/plain")
		res.Body("Whoopsie, my bad\n")
	default:
		res.SetStatusCode(response.StatusCode_200)
		res.SetDefaultHeaders()
		res.Headers.Set("Content-Type", "text/plain")
		res.Body("All good frfr\n")
	}
}
