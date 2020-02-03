package main

import (
	"log"
	"net"
)

var ladd net.TCPAddr
var radd net.TCPAddr

func main() {
	ladd.IP = net.ParseIP("localhost")
	ladd.Port = 8088

	radd.IP = net.ParseIP("localhost")
	radd.Port = 3000

	logConn, err := net.DialTCP("tcp", &ladd, &radd)
	if err != nil {
		log.Fatalf("Unable to connect logger to host: %+v", err)
	}

	// defer close
	logConn.Close()

	for {
		buffer := make([]byte, 1024)
		lenBuf, err := logConn.Read(buffer)
		if err != nil {
			log.Fatalf("Unable to read from buffer: %+v", err)
		}
		for i := 0; i < lenBuf; i++ {
			// check for null (\u0000) and remove trailing values
			if buffer[i] == '\u0000' {
				buffer = append(buffer[0:i])
				break
			}
			// if buffer not empty write buffer to log file
			if buffer != nil {
				// temp code to verify logic
				log.Print(string(buffer))
				// code to write log file
			}
		}
	}
}
