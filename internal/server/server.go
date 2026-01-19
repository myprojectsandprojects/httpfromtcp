package server

import (
	"boot.theprimeagen.tv/internal/request"
	"boot.theprimeagen.tv/internal/response"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener *net.Listener //@ pointer?
	closed   atomic.Bool
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	Code    response.StatusCode
	Message string
}

func (s *Server) listen(connHandler Handler) {
	// It mostly doesn't run when the server is closed because this function is waiting on Accept() most of the time
	//@ Is this fine?
	defer func() {
		(*s.listener).Close()
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

		conn, err := (*s.listener).Accept()
		if err != nil {
			log.Fatal(err) //@ How to report this error?
		}
		fmt.Print("Accepted a connection.\n")

		go s.handle(conn, connHandler)
	}
}

func (s *Server) handle(conn net.Conn, connHandler Handler) {
	defer func() {
		conn.Close()
		fmt.Print("Closed the connection.\n")
	}()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Printf("%v %v\n", req.RequestLine.Method, req.RequestLine.RequestTarget)
	// for k, v := range req.Headers {
	// 	fmt.Printf("%v: %v\n", k, v)
	// }

	var body bytes.Buffer
	handlerErr := connHandler(&body, req)
	if handlerErr != nil {
		err = response.WriteStatusLine(conn, handlerErr.Code)
		body.Write([]byte(handlerErr.Message))
	} else {
		err = response.WriteStatusLine(conn, response.StatusCode_200)
	}
	if err != nil {
		panic("boom")
	}

	contentLen := body.Len()

	headers := response.GetDefaultHeaders(contentLen)
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Fatal(err) //@ How to report this error?
	}

	_, err = conn.Write([]byte("\r\n"))
	if err != nil {
		log.Fatal(err) //@ How to report this error?
	}

	_, err = body.WriteTo(conn)
	if err != nil {
		log.Fatal(err) //@ How to report this error?
	}
}

func (s *Server) Close() error {
	log.Println("Server.Close() called")
	s.closed.Store(true)
	return nil
}

func (s *Server) Addr() net.Addr {
	return (*s.listener).Addr()
}

func Serve(port int, connHandler Handler) (*Server, error) {
	addr := fmt.Sprintf("127.0.0.1:%v", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &Server{listener: &l} //@ pointer?
	// s.closed.Store(false)      // seems to be "false" by default

	go s.listen(connHandler)

	return s, nil
}
