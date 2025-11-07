package main

import (
	"fmt"
	"os"
	"log"
	"io"
	"bytes"
)

func main() {
	f, err := os.Open("message.txt")
	if err != nil {
		log.Fatal(err)
	}
	// err = f.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	b := make([]byte, 8)
	s := ""
	for {
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
			fmt.Printf("read line: %#v\n", s)
			s = ""
			if i < len(data_read)-1 {
				s += string(data_read[i+1:])
			}
		} else {
			s += string(data_read)
		}
	}
	if len(s) > 0 {
		fmt.Printf("read last line: %#v\n", s)
	}
}
