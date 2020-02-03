// Package tcprproxy is a TCP reverse proxy adaptation
package tcprproxy

import (
	"fmt"
	"log"
	"net"
	"time"
)

var connSig, incoming, outgoing, done, exit chan string
var connections []net.Conn

// TCPForward uses tcp connections to establish a reverse proxy
func TCPForward(front string, back string) {
	listn, err := net.Listen("tcp", front) // create tcp connection
	if err != nil {
		log.Fatalf("failed to setup listener %v", err)
	}
	log.Println("ReverseProxy Listening on " + front)

	connSig = make(chan string)

	for {
		go session(listn, back)
		fmt.Println((<-connSig))
	}
}

func session(listn net.Listener, back string) {
	// wait for connection
	frontConn, err := listn.Accept()
	if err != nil {
		log.Fatalf("Failed to accept connection %v", err)
	}
	log.Println("Accepted Connection")

	// defer close
	defer frontConn.Close()

	connections = append(connections, frontConn)
	// send message to allow for an new connection
	connSig <- "Connection Complete"

	// create connection for server end
	serverConn, err := net.Dial("tcp", back)
	if err != nil {
		msg := "Dial failed for address" + back
		log.Fatalf("%+v %+v", msg, err)
	}

	// defer close
	defer serverConn.Close()

	// create message channels
	done = make(chan string)
	incoming = make(chan string)
	outgoing = make(chan string)

	// listen for message from client and log request
	go tcpListen(frontConn, incoming)
	// serve message from client to server
	go tcpServe(serverConn, incoming)
	// listen for response from server and log request
	go tcpListen(serverConn, outgoing)
	// server message from server to client
	go tcpServe(frontConn, outgoing)
	log.Println(<-done)
}

func tcpListen(conn net.Conn, packet chan string) {
	for {
		rBuf := make([]byte, 1024)

		err := conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			log.Fatalf("Unable to set connection (%s) deadline %v", conn.LocalAddr().String(), err)
		}

		_, err = conn.Read(rBuf)
		if err != nil {
			log.Printf("Error reading packet form %s: %+v", conn.LocalAddr().String(), err)
			break
		}
		// log.Println("Sent request ([]byte)::", requestBuf)
		log.Println("Got request/response (String)::\n", string(rBuf))

		packet <- string(rBuf)
	}
	done <- "Session Ended"
}

func tcpServe(conn net.Conn, packet chan string) {
	for {
		mBuf := <-packet

		//log.Println("Sent request/response (String)\n\n", string(mBuf))

		conn.Write([]byte(mBuf))
	}
}
