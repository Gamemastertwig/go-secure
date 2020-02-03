// Package httprproxy is encapsilation of the httputil.ReverseProxy
package httprproxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// HTTPForward uses httputil.ReverseProxy to 'forward' a request
// from client (in) to a server (out)
func HTTPForward(in string, out string) {
	// parsing destination url
	fmt.Println("Parsing Destination URL...")
	out = "http://" + out
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
