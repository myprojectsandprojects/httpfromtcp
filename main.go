package main

import (
	"fmt"
	"os"
	"log"
	"io"
)

func main() {
	f, err := os.Open("message.txt")
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 8)
	for {
		n, err := f.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Printf("read %v bytes: %s\n", n, b[:n])
	}
}
