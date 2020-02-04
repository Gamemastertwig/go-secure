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
func HTTPForward(in string, out string, logger string) {
	// attempt to connect logging server
	l.ConnectLogger(logger, logConn)

	// parsing destination url
	log.Println("HTTP Reverse Proxy")
	l.LogMessage("LOG", "HTTP Reverse Proxy", logConn)

	// need to add 'http://' to parse url for http reverse proxy
	out = "http://" + out
	to, err := url.Parse(out)
	if err != nil {
		l.LogMessage("ERROR", "Failed to parse "+out+":: "+err.Error(), logConn)
		log.Fatalln("Failed to parse " + out + ":: " + err.Error())
	}
	log.Println("Parsing Destination URL")
	l.LogMessage("LOG", "Parsing Destination URL", logConn)

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

	err = http.ListenAndServe(in, nil)
	if err != nil {
		l.LogMessage("ERROR", "Failed to start http.ListenAndServe:: "+err.Error(), logConn)
		log.Fatalln("Failed to start http.ListenAndServe:: " + err.Error())
	}
}
