package server

import (
	"boot.theprimeagen.tv/internal/response"
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

		conn, err := (*s.listener).Accept()
		if err != nil {
			log.Fatal(err) //@ How to report this error?
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	body := "Oh yeah, baby!"
	contentLen := len(body)

	err := response.WriteStatusLine(conn, response.StatusCode_500)
	// err := response.WriteStatusLine(conn, response.StatusCode_200)
	if err != nil {
		log.Fatal(err) //@ How to report this error?
	}

	headers := response.GetDefaultHeaders(contentLen)
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Fatal(err) //@ How to report this error?
	}

	_, err = conn.Write([]byte("\r\n"))
	if err != nil {
		log.Fatal(err) //@ How to report this error?
	}

	_, err = conn.Write([]byte(body))
	if err != nil {
		log.Fatal(err) //@ How to report this error?
	}

	conn.Close()
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
