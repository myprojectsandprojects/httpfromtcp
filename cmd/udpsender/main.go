package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(addr)
	fmt.Printf("IP: %v, Port: %v, Zone: %s\n", addr.IP, addr.Port, addr.Zone)

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println(conn)

	reader := bufio.NewReader(os.Stdin)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("got EOF")
				break
			}
			log.Fatal(err)
		}

		fmt.Println("user input:", input)
		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Fatal(err)
		}
	}
}
