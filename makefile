revp: 
	go run ./rproxy/rproxy.go
load:
	go run ./ldbal/ldbal.go
log:
	go run ./logger/logger.go
wall:
	go run ./fwall/fwall.go
intrus:
	go run ./ids/ids.go

all: revp load log wall intrus