// Package loghelper is a package to assist other applications
// with connecting o the logging server
package loghelper

// imports
import (
	"log"
	"net"
	"time"
)

// ConnLogMess creates a connection the log server and sends it a message
// then closes it
func ConnLogMess(logAddr string, mType string, message string) {
	var logConn net.Conn
	currentTime := time.Now()

	// create connection to logging server
	logConn, _ = net.Dial("tcp", logAddr)
	if logConn == nil {
		///log.Printf("Dial to logger at %+v failed:: %+v ", logAddr, err)
	} else {
		log.Println("Connected to logging server at " + logConn.LocalAddr().String())
		mes := currentTime.Format("2006/01/02 15:04:05") + " " + mType + " " + message
		logConn.Write([]byte(mes))
	}
}
