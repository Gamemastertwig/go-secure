// Package loghelper is a package to assist other applications
// with connecting o the logging server
package loghelper

// imports
import (
	"log"
	"net"
)

func ConnLogMess(logAddr string, mType string, message string) {
	var logConn net.Conn

	// create connection to logging server
	logConn, _ = net.Dial("tcp", logAddr)
	if logConn == nil {
		///log.Printf("Dial to logger at %+v failed:: %+v ", logAddr, err)
	} else {
		log.Println("Connected to logging server at " + logConn.LocalAddr().String())
		mes := mType + ": " + logConn.LocalAddr().String() + " " + message
		logConn.Write([]byte(mes))
	}
}
