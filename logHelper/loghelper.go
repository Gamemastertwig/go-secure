// Package loghelper is a package to assist other applications
// with connecting o the logging server
package loghelper

// imports
import (
	"errors"
	"log"
	"net"
)

// ConnectLogger attempts to onnect to logging sserver at logAddr (string)
func ConnectLogger(logAddr string, logConn net.Conn) {
	// IF logger (string) NOT "" THEN attempt connection with
	// log server
	if logAddr != "" {
		// create connection to logging server
		logConn, err := net.Dial("tcp", logAddr)
		if logConn == nil {
			log.Printf("Dial to logger at %+v failed:: %+v ", logAddr, err)
		}
		log.Println("Connected to logging server at " + logConn.LocalAddr().String())
		LogMessage("LOG", "Connected to logging server", logConn)
	}
}

// LogMessage send the message to the log server, it requires a message type (mType string)
// and a message (message string) using log server connectio (net.Conn) if fails returns
// an error
func LogMessage(mType string, message string, logConn net.Conn) error {
	if logConn != nil {
		mes := mType + ": " + logConn.LocalAddr().String() + " " + message
		logConn.Write([]byte(mes))
		return nil
	}
	log.Println("No connection to log server found: " + message)
	return errors.New("No connection to log server found")
}
