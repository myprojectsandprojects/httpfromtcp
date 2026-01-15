package main

import (
	"boot.theprimeagen.tv/internal/response"
	"boot.theprimeagen.tv/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	// "time"
)

// const port = 1

const port = 42069

func main() {
	err := response.WriteStatusLine(os.Stdout, response.StatusCode_200)
	if err != nil {
		log.Print(err)
	}

	headers := response.GetDefaultHeaders(3)
	err = response.WriteHeaders(os.Stdout, headers)
	if err != nil {
		log.Print(err)
	}

	server, err := server.Serve(port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started at address:", server.Addr().String())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")

	// time.Sleep(time.Second * 3)
}
