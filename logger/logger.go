package main

import (
	"bytes"
	"fmt"
	"net"
)

var connSig chan string

func main() {
	connSig = make(chan string)
	port := "localhost:9090"

	listn, _ := net.Listen("tcp", port)

	fmt.Println("Logging Server listening on " + port)

	for {
		go logSession(listn)
		<-connSig
	}
}

func logSession(listn net.Listener) {
	conn, _ := listn.Accept()
	defer conn.Close()

	fmt.Println("New Connection On " + conn.LocalAddr().String())
	connSig <- "Done"

	for {
		buffer := make([]byte, 1024)

		// Attempt read
		_, err := conn.Read(buffer)
		if err != nil {
			break
		}

		cleanBuf := bytes.Trim(buffer, "\x00")

		// display log information here
		// just display to STDOUT for now
		if string(cleanBuf) != "" {
			fmt.Println(string(cleanBuf))
		}
	}
}
