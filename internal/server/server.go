package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener *net.Listener //@ pointer?
	closed   atomic.Bool
}

func (s *Server) listen() {
	// It mostly doesn't run when the server is closed because this function is waiting on Accept() most of the time
	//@ Is this fine?
	// defer func() {
	// 	(*s.listener).Close()
	// 	log.Print("closed the listener")
	// }()

	for {
		if s.closed.Load() {
			// It mostly doesn't run when the server is closed because this function is waiting on Accept() most of the time
			//@ Is this fine?
			// Accept no more connections when .Close() is called on the server
			log.Println("server done serving")
			break
		}

		c, err := (*s.listener).Accept()
		if err != nil {
			log.Fatalln("fail")
		}

		go s.handle(c)
	}
}

func (s *Server) handle(c net.Conn) {
	body := "Hello world!"
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(body), body)
	// response := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: 21\r\n\r\n<h1>Hello World!</h1>")
	// response := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 21\r\n\r\n<h1>Hello World!</h1>")
	n, err := c.Write([]byte(response))
	if err != nil {
		log.Fatalln("fail")
	}
	log.Printf("Wrote a response (%v bytes)", n)

	c.Close()
}

func (s *Server) Close() error {
	log.Println("Server.Close() called")
	s.closed.Store(true)
	return nil
}

func (s *Server) Addr() net.Addr {
	return (*s.listener).Addr()
}

func Serve(port int) (*Server, error) {
	addr := fmt.Sprintf("127.0.0.1:%v", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	s := &Server{listener: &l} //@ pointer?
	s.closed.Store(false)      // seems to be "false" by default
	go s.listen()
	return s, nil
}
