package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var front, back string
var connSig, incoming, outgoing, done, exit chan string
var connections []net.Conn

func main() {
	// front = "localhost:8080"
	// back = "http://localhost:3000/"

	// httpForward(front, back)

	// TCP reverse proxy is failing, need to find
	// way to restablish a connection after EOF

	front = "localhost:8080"
	back = "localhost:3000"

	tcpForward()

}

func httpForward(in string, out string) {
	// parsing destination url
	fmt.Println("Parsing Destination URL...")
	to, err := url.Parse(out)
	if err != nil {
		log.Fatalln(err)
	}

	// modify http header
	fmt.Println("Modifing header...")
	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", to.Host)
		req.URL.Scheme = "http"
		req.URL.Host = to.Host
	}

	// prepare to forward
	fmt.Println("Preparing to forward...")
	proxy := &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
	http.HandleFunc("/detail", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(in, nil))
}

// not working correctly at the moment
func tcpForward() {
	listn, err := net.Listen("tcp", front) // create tcp connection
	if err != nil {
		log.Fatalf("failed to setup listener %v", err)
	}
	log.Println("ReverseProxy Listening on " + front)

	connSig = make(chan string)

	for {
		go session(listn)
		fmt.Println((<-connSig))
	}
}

func session(listn net.Listener) {
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
