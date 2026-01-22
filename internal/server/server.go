package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"unsafe"

	"boot.theprimeagen.tv/internal/request"
	"boot.theprimeagen.tv/internal/response"
)

type Server struct {
	Listener net.Listener //@ pointer?
	closed   atomic.Bool
}

type Handler func(res *response.Response, req *request.Request)

type HandlerError struct {
	Code    response.StatusCode
	Message string
}

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

	res := response.New()
	handler(res, req)

	res.Bytes.WriteTo(conn)
}

func (s *Server) Close() error {
	log.Println("Server.Close() called")
	s.closed.Store(true)
	return nil
}

func (s *Server) Addr() net.Addr {
	return s.Listener.Addr()
}

func Serve(port int, handler Handler) (*Server, error) {
	addr := fmt.Sprintf("127.0.0.1:%v", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	fmt.Printf("listener's type is %T\n", l)
	fmt.Printf("listener's size is %v\n", unsafe.Sizeof(l))

	s := &Server{Listener: l} //@ pointer?
	// s.closed.Store(false)      // seems to be "false" by default

	go s.listen(handler)

	return s, nil
}
