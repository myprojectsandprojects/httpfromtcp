package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"boot.theprimeagen.tv/internal/request"
	"boot.theprimeagen.tv/internal/response"
)

type Server struct {
	Listener net.Listener
	closed   atomic.Bool
}

type Handler func(w *response.Writer, req *request.Request)

func (s *Server) listen(handler Handler) {
	// It mostly doesn't run when the server is closed because this function is waiting on Accept() most of the time
	//@ Is this fine?
	defer func() {
		s.Listener.Close()
		log.Print("closed the listener")
	}()

	for {
		if s.closed.Load() {
			// It mostly doesn't run when the server is closed because this function is waiting on Accept() most of the time
			//@ Is this fine?
			// Accept no more connections when .Close() is called on the server
			log.Println("server done serving")
			break
		}

		conn, err := s.Listener.Accept()
		if err != nil {
			log.Fatal(err) //@ How to report this error?
		}

		go s.handle(conn, handler)
	}
}

func (s *Server) handle(conn net.Conn, handler Handler) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Printf("%v %v\n", req.RequestLine.Method, req.RequestLine.RequestTarget)

	w := response.NewWriter(conn)
	handler(&w, req)
}

func (s *Server) Close() error {
	log.Println("Server.Close() called")
	s.closed.Store(true)
	return nil
}

func Serve(port int, handler Handler) (*Server, error) {
	addr := fmt.Sprintf("127.0.0.1:%v", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &Server{Listener: l}

	go s.listen(handler)

	return s, nil
}
