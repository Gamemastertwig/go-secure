// Package tcprproxy is a TCP reverse proxy adaptation
package tcprproxy

// Import packages
import (
	"log"
	"math/rand"
	"net"
	"time"

	l "github.com/Gamemastertwig/go-secure/logHelper"
)

// Global Variables
var connSig, incoming, outgoing, done, exit chan string
var logConn net.Conn

// TCPForward uses tcp connections to establish a reverse proxy
func TCPForward(front string, back string, logger string) {
	// create tcp connection
	listn, err := net.Listen("tcp", front)
	if err != nil {
		l.ConnLogMess(logger, "ERROR RPROXY: ", front+" Failed to setup listiner for RPROXY")
		log.Fatalf("Failed to setup listener %v", err)
	}

	l.ConnLogMess(logger, "LOG RPROXY: ", "(TCP) Listening on "+front)
	log.Println("TCP ReverseProxy Listening on " + front)

	// create connection signal channel
	connSig = make(chan string)

	// allow for multiple connections
	for {
		go Session(listn, back, logger)
		<-connSig
	}
}

// TCPForwardLB uses tcp connections to establish a reverse proxy with load balancer
func TCPForwardLB(front string, backends []string, logAddr string) {
	// create tcp connection
	listn, err := net.Listen("tcp", front)
	if err != nil {
		l.ConnLogMess(logAddr, "ERROR LB: ", front+" Failed to setup listiner for LB RPROXY")
		log.Fatalf("Failed to setup listener %v", err)
	}

	l.ConnLogMess(logAddr, "LOG LB: ", "(TCP) Listening on "+front)
	log.Println("Load Balancer Listening on " + front)

	// create connection signal channel
	connSig = make(chan string)

	// allow for multiple connections
	for {
		// call loadBalance here
		bAddr := loadBalanceRand(backends, logAddr)
		go Session(listn, bAddr, logAddr)
		<-connSig
	}
}

func loadBalanceRand(backends []string, logAddr string) string {
	// verify back ends
	for i, b := range backends {
		TempConn, err := net.Dial("tcp", b)
		if err != nil {
			// remove from slice ([]sting) if it fails to make a connection
			backends = append(backends[:i], backends[i+1:]...)
		}
		TempConn.Close()
	}

	// seed rand based on time per OS
	rand.Seed(time.Now().UnixNano())
	// get a random int from 0 to length of string slice -1
	lenBack := len(backends)
	n := rand.Intn(lenBack)

	l.ConnLogMess(logAddr, "LOG LB: ", "Sending load to "+backends[n])

	return backends[n]
}

// Session allows for multiple connections from clients at the same time
// listening on front end (net.Listener) and then connects to back end
// address (backAddr string)
func Session(listn net.Listener, backAddr string, logAddr string) {
	// wait for front end connection
	frontConn, err := listn.Accept()
	if err != nil {
		l.ConnLogMess(logAddr, "ERROR SESSION: ", "Failed to accept connection:: "+err.Error())
		log.Fatalln("Failed to accept connection:: " + err.Error())
	}
	l.ConnLogMess(logAddr, "LOG SESSION: ", "Accepted Connection from "+frontConn.LocalAddr().String())

	// defer close : member LIFO
	defer frontConn.Close()

	// send message to allow for an new connection
	connSig <- "Done"

	// create connection for server end
	serverConn, err := net.Dial("tcp", backAddr)
	if err != nil {
		l.ConnLogMess(logAddr, "ERROR SESSION: ", "Dial failed for address "+backAddr+":: "+err.Error())
		log.Fatalln("Dial failed for address " + backAddr + ":: " + err.Error())
	}
	l.ConnLogMess(logAddr, "LOG SESSION: ", "Dial succesful to "+backAddr)

	// defer close
	defer serverConn.Close()

	// create message channels
	done = make(chan string)
	incoming = make(chan string)
	outgoing = make(chan string)

	// listen for message from client and log request
	go TCPListen(frontConn, incoming, logAddr)
	// serve message from client to server
	go TCPServe(serverConn, incoming, logAddr)
	// listen for response from server and log request
	go TCPListen(serverConn, outgoing, logAddr)
	// server message from server to client
	go TCPServe(frontConn, outgoing, logAddr)
	<-done
}

// TCPListen listens on conn (net.Conn) and reads the data ([]byte) into a
// string channel
func TCPListen(conn net.Conn, packet chan string, logAddr string) {
	for {
		// create read buffer
		rBuf := make([]byte, 1024)

		err := conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			l.ConnLogMess(logAddr, "ERROR TCPLISTEN: ", "Failed deadline setup for "+
				conn.LocalAddr().String()+":: "+err.Error())
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
		l.ConnLogMess(logAddr, "LOG TCPLISTEN: ", "Packet read.")
	}
	done <- "Done"
}

// TCPServe serves the data ([]byte) from the packet (string channel)
// to the connection (net.Conn)
func TCPServe(conn net.Conn, packet chan string, logAddr string) {
	for {
		// get buffer from packet channel
		mBuf := <-packet

		// connection exist write buffer to connection
		if conn != nil {
			conn.Write([]byte(mBuf))
			l.ConnLogMess(logAddr, "LOG TCPSERVE: ", "Packet sent.")
		}
	}
}
