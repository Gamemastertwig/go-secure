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
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"strconv"
)

// config/rule struct
type config struct {
	Proto string   `json:"proto,omitempty"`
	Port  int      `json:"port,omitempty"`
	Allow []string `json:"allow,omitempty"`
}

func init() {
	// parse for arg
	flag.Parse()
}

func main() {
	// get path of config (json) from arg
	path := flag.Arg(0)
	if path == "" {
		log.Fatal("Did you set a config file?")
	}

	// open config file (json) at path
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Unable to open config: %+v", err)
	}
	// defer close
	defer f.Close()

	// decode config (json)
	var configs []config
	err = json.NewDecoder(f).Decode(&configs)
	if err != nil {
		log.Fatalf("Unable to decode config: %+v", err)
	}

	// for each config in config (json) file process rule
	for _, c := range configs {
		port := strconv.Itoa(c.Port) // convert port to string to use as arg
		// Append rule to chain
		cmd := exec.Command("iptables", "-A", "INPUT", "-p", c.Proto, "--dport", port, "-j", "REJECT")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Unable to append rule to chain: %+v", err)
		}
		log.Print(string(output))

		// Allow each ip in the 'Allow' slice
		for _, a := range c.Allow {
			cmd := exec.Command("iptables", "-I", "INPUT", "-s", a, "-p", c.Proto, "--dport", port, "-j", "ACCEPT")
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatalf("Unable to insert allowed ip per rule: %+v", err)
			}
			log.Print(string(output))
		}
	}
}
