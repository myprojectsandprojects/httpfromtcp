package main

import (
	"boot.theprimeagen.tv/internal/request"
	"fmt"
	"log"
	"net"
	// "os"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		log.Fatal(err)
	}

	addr := l.Addr()
	fmt.Printf("listening on %v, %v\n", addr.Network(), addr.String())

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("connection accepted:", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %v\n", req.RequestLine.Method)
		fmt.Printf("- Target: %v\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %v\n", req.RequestLine.HttpVersion)

		fmt.Printf("Headers:\n")
		for k, v := range req.Headers {
			fmt.Printf("- %v: %v\n", k, v)
		}

		if len(req.Body) > 0 {
			fmt.Printf("Body: %q\n", req.Body)
		}

		conn.Close()
	}
}
