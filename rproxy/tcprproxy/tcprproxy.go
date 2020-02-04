// Package tcprproxy is a TCP reverse proxy adaptation
package tcprproxy

// Import packages
import (
	"log"
	"net"
	"time"

	l "github.com/Gamemastertwig/go-secure/logHelper"
)

// Global Variables
var connSig, incoming, outgoing, done, exit chan string
var logConn net.Conn

// TCPForward uses tcp connections to establish a reverse proxy
func TCPForward(front string, back string, logger string) {
	// attempt to connect logging server
	l.ConnectLogger(logger, logConn)

	// create tcp connection
	listn, err := net.Listen("tcp", front)
	if err != nil {
		temp := front + ": Failed to setup listiner: "
		l.LogMessage("ERROR", temp, logConn)
		log.Fatalf("Failed to setup listener %v", err)
	}
	l.LogMessage("LOG", "ReverseProxy Listening on "+front, logConn)
	log.Println("TCP ReverseProxy Listening on " + front)

	// create connection signal channel
	connSig = make(chan string)

	// allow for multiple connections
	for {
		go Session(listn, back)
		<-connSig
	}
}

// Session allows for multiple connections from clients at the same time
// listening on front end (net.Listener) and then connects to back end
// address (backAddr string)
func Session(listn net.Listener, backAddr string) {
	// wait for front end connection
	frontConn, err := listn.Accept()
	if err != nil {
		l.LogMessage("ERROR", "Failed to accept connection:: "+err.Error(), logConn)
		log.Fatalln("Failed to accept connection:: " + err.Error())
	}
	l.LogMessage("LOG", "Accepted Connection from"+frontConn.LocalAddr().String(), logConn)

	// defer close : member LIFO
	defer frontConn.Close()
	defer l.LogMessage("LOG", "Closing connection: "+frontConn.LocalAddr().String(), logConn)

	// send message to allow for an new connection
	connSig <- "Done"

	// create connection for server end
	serverConn, err := net.Dial("tcp", backAddr)
	if err != nil {
		l.LogMessage("ERROR", "Dial failed for address "+backAddr+":: "+err.Error(), logConn)
		log.Fatalln("Dial failed for address " + backAddr + ":: " + err.Error())
	}
	l.LogMessage("LOG", "Dial succesful to "+frontConn.LocalAddr().String(), logConn)

	// defer close
	defer serverConn.Close()
	defer l.LogMessage("LOG", "Closing connection: "+serverConn.LocalAddr().String(), logConn)

	// create message channels
	done = make(chan string)
	incoming = make(chan string)
	outgoing = make(chan string)

	// listen for message from client and log request
	go TCPListen(frontConn, incoming)
	// serve message from client to server
	go TCPServe(serverConn, incoming)
	// listen for response from server and log request
	go TCPListen(serverConn, outgoing)
	// server message from server to client
	go TCPServe(frontConn, outgoing)
	<-done
}

// TCPListen listens on conn (net.Conn) and reads the data ([]byte) into a
// string channel
func TCPListen(conn net.Conn, packet chan string) {
	for {
		// create read buffer
		rBuf := make([]byte, 1024)

		err := conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			l.LogMessage("ERROR", "Failed deadline setup for "+
				conn.LocalAddr().String()+":: "+err.Error(), logConn)
			log.Fatalln("Failed deadline setup for " +
				conn.LocalAddr().String() + ":: " + err.Error())
		}

		// Attempt read
		_, err = conn.Read(rBuf)
		if err != nil {
			break
		}

		// place buffer in packet channel
		packet <- string(rBuf)
	}
	done <- "Done"
}

// TCPServe serves the data ([]byte) from the packet (string channel)
// to the connection (net.Conn)
func TCPServe(conn net.Conn, packet chan string) {
	for {
		// get buffer from packet channel
		mBuf := <-packet

		// connection exist write buffer to connection
		if conn != nil {
			conn.Write([]byte(mBuf))
		}
	}
}
