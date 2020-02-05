// Package launcher is a helper package for starting or stoping launchable application in go
package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Program is a struct to hold variables for a launchable application
type Program struct {
	Dir  string
	Comm string
	Arg  string
	Pid  int
}

// StartStop checks if program is running and calls Start or Stop
func StartStop(cmd Program) {

	// check if process running cmd.pid == 0
	if cmd.Pid == 0 {
		num, pass := Check(cmd)
		if pass {
			cmd.Pid = num
		}
	}

	// check cmd.pid again
	if cmd.Pid == 0 {
		// start program
		Start(cmd)
	} else {
		// stop program
		Stop(cmd)
	}
}

// Start will start a launchable application
func Start(cmd Program) {
	err := os.Chdir(cmd.Dir)
	if err != nil {
		fmt.Println(err)
	}

	proc := exec.Command(cmd.Comm)
	fmt.Println("Starting " + cmd.Comm)
	proc.Start()
}

// Stop will kill (want to change this later) a launchable application
func Stop(cmd Program) {
	proc := exec.Command("kill", "-9", strconv.Itoa(cmd.Pid))
	fmt.Println("Killing " + cmd.Comm)
	proc.Run()
}

// Check looks to see if launchable application (process) is running.
// Returns the PID (int) and TRUE (bool) if found.
// Returns 0 (int) and FALSE (bool) if not found.
func Check(cmd Program) (int, bool) {
	var progName string
	pID := 0

	// removes "./" if present in cmd.comm
	if strings.Contains(cmd.Comm, "./") {
		progName = strings.ReplaceAll(cmd.Comm, "./", "")
	} else {
		progName = cmd.Comm
	}

	// get all processes running
	check := exec.Command("ps", "-a")
	// place output into buffer
	buffer, err := check.Output()
	if err != nil {
		fmt.Println("Error::" + err.Error())
	}
	// split buffer at new line
	temp := strings.Split(string(buffer), "\n")
	// search each line (process) for program name
	for _, s := range temp {
		if strings.Contains(s, progName) {
			temp2 := strings.Split(s, " ")
			if temp2[0] != "" {
				pID, err = strconv.Atoi(temp2[0])
				if err != nil {
					fmt.Println("Failed temp2[0]: " + err.Error())
				}
			} else {
				pID, err = strconv.Atoi(temp2[1])
				if err != nil {
					fmt.Println("Failed temp2[0]: " + err.Error())
				}
			}
		}
	}
	if pID == 0 {
		return pID, false
	}
	return pID, true
}
