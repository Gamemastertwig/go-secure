package logwriter

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func CheckForFile(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateFile(filename string) bool {
	_, err := os.Create(filename)
	if err != nil {
		return false
	}
	return true
}

func WriteNew(filename string, message []byte) bool {
	if !CheckForFile(filename) {
		if !CreateFile(filename) {
			log.Println("Unable to find or create file: " + filename)
			return false
		}
	} else {
		f, err := os.Open(filename)
		if err != nil {
			log.Println("Unable to open file: " + filename)
			return false
		}
		defer f.Close()

		fmt.Fprint(f, string(message))
	}
	return true
}

func WriteAppend(filename string, message []byte) bool {
	if !CheckForFile(filename) {
		if !CreateFile(filename) {
			log.Println("Unable to find or create file: " + filename)
			return false
		}
	} else {
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			log.Println("Unable to open file: " + filename)
			return false
		}
		defer f.Close()

		fmt.Fprintln(f, string(message))
	}
	return true
}

func WriteRaw(filename string, message []byte) bool {
	if !CheckForFile(filename) {
		if !CreateFile(filename) {
			log.Println("Unable to find or create file: " + filename)
			return false
		}
	} else {
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			log.Println("Unable to open file: " + filename)
			return false
		}
		defer f.Close()

		fmt.Fprint(f, message)
	}
	return true
}

func ReadFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("File reading error", err)
		return nil
	}
	return data
}
