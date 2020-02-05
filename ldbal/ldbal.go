package main

import "github.com/Gamemastertwig/go-secure/rproxy/tcprproxy"

func main() {

	front := "localhost:8081"
	backends := []string{"localhost:3000", "localhost:3001", "localhost:3002"}
	logger := "localhost:9090"

	tcprproxy.TCPForwardLB(front, backends, logger)
}
