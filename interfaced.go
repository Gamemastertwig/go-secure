package main

import "github.com/Gamemastertwig/go-secure/launcher"

func main() {
	logger := launcher.Program{Dir: "./logger", Comm: "./logger", Arg: "", Pid: 0}
	proxy := launcher.Program{Dir: "./rproxy", Comm: "./rproxy", Arg: "", Pid: 0}

	_, started := launcher.Check(logger)
	if !started {
		launcher.Start(logger)
	}

	_, started = launcher.Check(proxy)
	if !started {
		launcher.Start(proxy)
	}
}
