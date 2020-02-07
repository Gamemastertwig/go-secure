// GO-SECURE (FWALL)
/*
	GO-SECURE (FWALL) is part of a suite of security applications built in Go.
	This modual is a light-weight tool for modifying the iptables to create
	firewall (FWALL) rules. This modual requires super user acess (sudo) and
	must be ran on the server (hardware, vm, docker, etc) as the server
	applications it is meant to protect. (Academic concept)

	Revature: Brandon Locker (GameMasterTwig)
*/
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/Gamemastertwig/go-secure/loghelper"
	"github.com/Gamemastertwig/go-secure/logwriter"
)

// config/rule struct
type config struct {
	Proto string   `json:"proto"`
	Port  int      `json:"port"`
	Allow []string `json:"allow"`
}

// connection to logger server
type logger struct {
	LogAddress string `json:"logger"`
}

var configs []config
var logs []logger
var logAddr string

func init() {
	flag.Parse()
	getLogAdder()
}

func getLogAdder() {
	if logwriter.CheckForFile("logConfig.json") {
		// file is present
		f, err := ioutil.ReadFile("logConfig.json")
		if err != nil {
			log.Fatalf("Unable to open logConfig: %+v", err)
		}

		// decode config (json)
		err = json.Unmarshal(f, &logs)
		if err != nil {
			log.Fatalf("Unable to decode logConfig: %+v", err)
		}
		logAddr = logs[0].LogAddress
	}
}

func main() {
	temp := flag.Arg(0)
	fmt.Println(temp)
	switch temp {
	case "Add", "add":
		if flag.Arg(1) != "" {
			addRules(flag.Args())
		} else {
			fmt.Println("To add rules please enter as...\n" +
				"Add [protocal] [port] [allowed IPs (seperated by a space)]...\n" +
				"you can also modify the config.json file in this dir")
		}
		loghelper.ConnLogMess(logAddr, "FWALL ERROR:", "Failed to add rule to firewall")
	case "Display", "display":
		displayTables()
	case "Clear", "clear":
		clearRules()
	case "":
		applRules()
	}
}

func addRules(args []string) {
	readConfig("config.json")

	var temp config
	if args[1] != "" {
		temp.Proto = args[1]
	}
	if args[2] != "" {
		n, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
		}
		temp.Port = n
	}
	if args[3] != "" {
		allow := args[3:]
		for _, a := range allow {
			temp.Allow = append(temp.Allow, a)
		}
	}

	configs = append(configs, temp)

	writeConfig(configs, "config.json")

	applRules()
	loghelper.ConnLogMess(logAddr, "FWALL LOG:", "Added and applied new rule to firewall")
}

func clearRules() {
	cmd := exec.Command("iptables", "-F")
	_, err := cmd.CombinedOutput()
	if err != nil {
		loghelper.ConnLogMess(logAddr, "FWALL ERROR:", "Unable to FLUSH iptables: "+err.Error())
		log.Fatalf("Unable to FLUSH iptables: %+v", err)
	}
	loghelper.ConnLogMess(logAddr, "FWALL LOG:", "Flushed all rules from firewall (iptables)")
}

func readConfig(filename string) {
	// open config file (json) at path
	if logwriter.CheckForFile(filename) {
		// file is present
		f, err := os.Open(filename)
		if err != nil {
			loghelper.ConnLogMess(logAddr, "FWALL ERROR:", "Unable to open config: "+err.Error())
			log.Fatalf("Unable to open config: %+v", err)
		}
		// defer close
		defer f.Close()

		// decode config (json)
		err = json.NewDecoder(f).Decode(&configs)
		if err != nil {
			loghelper.ConnLogMess(logAddr, "FWALL ERROR:", "Unable to decode config: "+err.Error())
			log.Fatalf("Unable to decode config: %+v", err)
		}
	}
	loghelper.ConnLogMess(logAddr, "FWALL LOG:", "Read from config sucessful")
}

func writeConfig(data interface{}, filename string) (int, error) {
	//write data as buffer to json encoder
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		loghelper.ConnLogMess(logAddr, "FWALL ERROR:", "Unable to encode config: "+err.Error())
		return 0, err
	}
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		loghelper.ConnLogMess(logAddr, "FWALL ERROR:", "Unable to open config: "+err.Error())
		return 0, err
	}
	n, err := file.Write(buffer.Bytes())
	if err != nil {
		loghelper.ConnLogMess(logAddr, "FWALL ERROR:", "Unable to write to config: "+err.Error())
		return 0, err
	}
	loghelper.ConnLogMess(logAddr, "FWALL LOG:", "Write to config sucessful")
	return n, nil
}

func applRules() {
	// rules
	readConfig("config.json")

	if configs != nil {
		// for each config in config (json) file process rule
		for _, c := range configs {
			port := strconv.Itoa(c.Port) // convert port to string to use as arg
			// Append rule to chain
			cmd := exec.Command("iptables", "-A", "INPUT", "-p", c.Proto, "--dport", port, "-j", "REJECT")
			_, err := cmd.CombinedOutput()
			if err != nil {
				loghelper.ConnLogMess(logAddr, "FWALL ERROR:", "Unable to append rule to chain: "+err.Error())
				log.Fatalf("Unable to append rule to chain: %+v", err)
			}

			// Allow each ip in the 'Allow' slice
			for _, a := range c.Allow {
				cmd := exec.Command("iptables", "-I", "INPUT", "-s", a, "-p", c.Proto, "--dport", port, "-j", "ACCEPT")
				_, err := cmd.CombinedOutput()
				if err != nil {
					loghelper.ConnLogMess(logAddr, "FWALL ERROR:", "Unable to insert allowed ip per rule: "+err.Error())
					log.Fatalf("Unable to insert allowed ip per rule: %+v", err)
				}
			}
		}
		loghelper.ConnLogMess(logAddr, "FWALL LOG:", "Rules applied to firewall (iptables) sucessfully")
	} else {
		// file is not present
		loghelper.ConnLogMess(logAddr, "FWALL ERROR:", "No rules to apply!")
	}
}

func displayTables() {
	cmd := exec.Command("iptables", "-S")
	output, err := cmd.CombinedOutput()
	if err != nil {
		loghelper.ConnLogMess(logAddr, "FWALL ERROR:", "Unable to display iptables: "+err.Error())
		log.Fatalf("Unable to display iptables: %+v", err)
	}
	fmt.Println(string(output))
	loghelper.ConnLogMess(logAddr, "FWALL LOG:", "Displayed Firewall (iptables) sucessfully")
}
