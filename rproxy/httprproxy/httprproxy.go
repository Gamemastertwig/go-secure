// Package httprproxy is encapsilation of the httputil.ReverseProxy
package httprproxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	l "github.com/Gamemastertwig/go-secure/logHelper"
)

var logConn net.Conn

// HTTPForward uses httputil.ReverseProxy to 'forward' a request
// from client (in) to a server (out)
func HTTPForward(fAddr string, bAddr string, logAddr string) {
	// parsing destination url
	log.Println("HTTP Reverse Proxy")
	l.ConnLogMess(logAddr, "LOG", "HTTP Reverse Proxy")

	// need to add 'http://' to parse url for http reverse proxy
	bAddr = "http://" + bAddr
	to, err := url.Parse(bAddr)
	if err != nil {
		l.ConnLogMess(logAddr, "ERROR", "Failed to parse "+bAddr+":: "+err.Error())
		log.Fatalln("Failed to parse " + bAddr + ":: " + err.Error())
	}
	log.Println("Parsing Destination URL")
	l.ConnLogMess(logAddr, "LOG", "Parsing Destination URL")

	// modify http header
	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", to.Host)
		req.URL.Scheme = "http"
		req.URL.Host = to.Host
	}
	fmt.Println("Modifed header")

	// prepare to forward
	proxy := &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
	http.HandleFunc("/detail", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	err = http.ListenAndServe(fAddr, nil)
	if err != nil {
		l.ConnLogMess(logAddr, "ERROR", "Failed to start http.ListenAndServe:: "+err.Error())
		log.Fatalln("Failed to start http.ListenAndServe:: " + err.Error())
	}
}
