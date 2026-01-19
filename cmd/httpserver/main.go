package main

import (
	"boot.theprimeagen.tv/internal/request"
	"boot.theprimeagen.tv/internal/response"
	"boot.theprimeagen.tv/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	// "time"
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

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	var err *server.HandlerError = nil

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		err = &server.HandlerError{
			Code:    response.StatusCode_400,
			Message: "Your problem is not my problem\n",
		}
	case "/myproblem":
		err = &server.HandlerError{
			Code:    response.StatusCode_500,
			Message: "Whoopsie, my bad\n",
		}
	default:
		w.Write([]byte("All good frfr\n"))
		err = nil
	}

	return err
}
