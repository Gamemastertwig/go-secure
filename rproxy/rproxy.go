package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	forward("localhost:8080", "http://localhost:3000/")
}

func forward(in string, out string) {
	to, err := url.Parse(out)
	if err != nil {
		log.Fatalln(err)
	}

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", to.Host)
		req.URL.Scheme = "http"
		req.URL.Host = to.Host
	}

	proxy := &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(in, nil))
}
