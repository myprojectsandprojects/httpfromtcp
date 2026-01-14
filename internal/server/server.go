package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	closed atomic.Bool
}

func (s *Server) listen(port int) {
	addr := fmt.Sprintf("127.0.0.1:%v", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("error: net.Listen: %v", err)
	}
	defer l.Close()

	for {
		if s.closed.Load() {
			// Accept no more connections when .Close() is called on the server
			log.Println("server done serving")
			break
		}

		c, err := l.Accept()
		if err != nil {
			log.Fatalln("fail")
		}

		go s.handle(c)
	}
}

func (s *Server) handle(c net.Conn) {
	response := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 14\r\n\r\nHello World?!!")
	_, err := c.Write(response)
	if err != nil {
		log.Fatalln("fail")
	}

	c.Close()
}

func (s *Server) Close() error {
	log.Println("Server.Close() called")
	s.closed.Store(true)
	return nil
}

func Serve(port int) (*Server, error) {
	s := Server{}
	go s.listen(port)
	return &s, nil
}
