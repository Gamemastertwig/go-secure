package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"sync"
	"text/template"

	"github.com/Gamemastertwig/go-secure/loghelper"
)

// ContentTemplate holds data used to display content to the http template
type ContentTemplate struct {
	RproxyBody   string
	FirewallBody string
	LBalBody     string
	IDSBody      string
	LogBody      string
}

var fileData chan []byte

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

	loghelper.ConnLogMess("localhost:9090", "LOG FILE", ":9998")
	temp := listenForLogFile(":9998")
	//fmt.Print(temp)

	content.RproxyBody = "Reverse Proxy content get loaded here..."
	content.FirewallBody = "Firewall content get loaded here..."
	content.LBalBody = "[LB] Comming Soon..."
	content.IDSBody = "[IDS] Comming Soon..."
	content.LogBody = string(temp)

	css := http.FileServer(http.Dir("web/css/"))
	http.Handle("/css/", http.StripPrefix("/css/", css))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template.Execute(w, content)
	})
	fmt.Println("Attempting start")
	var wg sync.WaitGroup
	wg.Add(1)
	go http.ListenAndServe("localhost:9999", nil)
	wg.Wait()
}

func listenForLogFile(port string) []byte {
	var cleanBuf []byte

	if port == "" {
		port = ":9998"
	}

	listn, _ := net.Listen("tcp", port)

	fmt.Println("listening on " + port)

	conn, _ := listn.Accept()
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)

		// Attempt read
		_, err := conn.Read(buffer)
		if err != nil {
			break
		}

		cleanBuf := bytes.Trim(buffer, "\x00")

		if string(cleanBuf) != "" {
			fmt.Print(string(cleanBuf))
		}
	}
	return cleanBuf
}
