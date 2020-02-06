package main

import (
	"fmt"
	"net/http"
	"text/template"
)

type ContentTemplate struct {
	RproxyBody   string
	FirewallBody string
	LBalBody     string
	IDSBody      string
	LogBody      string
}

func main() {
	// logger := launcher.Program{Dir: "./logger", Comm: "./logger", Arg: "", Pid: 0}
	// proxy := launcher.Program{Dir: "./rproxy", Comm: "./rproxy", Arg: "", Pid: 0}

	// _, started := launcher.Check(logger)
	// if !started {
	// 	launcher.Start(logger)
	// }

	// _, started = launcher.Check(proxy)
	// if !started {
	// 	launcher.Start(proxy)
	// }
	var template = template.Must(template.ParseFiles("web/index.html"))
	var content ContentTemplate

	content.RproxyBody = "Reverse Proxy content get loaded here..."
	content.FirewallBody = "Firewall content get loaded here..."
	content.LBalBody = "[LB] Comming Soon..."
	content.IDSBody = "[IDS] Comming Soon..."
	content.LogBody = "Logger content gets loaded here..."

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template.Execute(w, content)
	})
	fmt.Println("Attempting start")

	err := http.ListenAndServe("localhost:9999", nil)
	if err != nil {
		fmt.Println(err)
	}
}
