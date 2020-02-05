// GO-SECURE (RPROXY)
/*
	GO-SECURE (RPROXY) is part of a suite of security applications built in Go.
	This modual is a tool for forwarding requests to a server (reverse proxy).
	You can currently forward using HTTP or TCP protocals.
	(Academic concept)

	Revature: Brandon Locker (GameMasterTwig)
*/
package main

import (
	h "github.com/Gamemastertwig/go-secure/rproxy/httprproxy"
	t "github.com/Gamemastertwig/go-secure/rproxy/tcprproxy"
)

func init() {

}

func main() {
	front := "localhost:8080"
	back := "localhost:3000"
	logger := "localhost:9090"
	useTCP := true

	if useTCP {
		t.TCPForward(front, back, logger)
	} else {
		h.HTTPForward(front, back, logger)
	}
}
