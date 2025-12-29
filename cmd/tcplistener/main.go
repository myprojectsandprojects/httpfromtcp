package main

import (
	"fmt"
	// "os"
	"bytes"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer f.Close() // can fail
		defer close(out)

		s := ""
		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			data_read := b[:n]

			// Multiple '\n' in one read are not handled correctly!!!
			if i := bytes.IndexByte(data_read, '\n'); i != -1 {
				s += string(data_read[:i])
				out <- s
				s = ""
				s += string(data_read[i+1:]) // low <= high <= capacity, no out of bounds error (if i is at the end, you get an empty slice)
			} else {
				s += string(data_read)
			}
		}
		if len(s) > 0 {
			out <- s
		}
	}()

	return out
}

func main() {
	// f, err := os.Open("message.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }

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
		fmt.Println("connection accepted:", conn)

		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}

		fmt.Println("connection closed")
	}
}
