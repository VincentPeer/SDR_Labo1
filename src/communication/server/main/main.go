package main

import (
	"SDR_Labo1/src/communication/server"
	"os"
)

var (
	connHost   = "localhost"
	connPort   = "3333"
	configPath = "./config.json"
)

func main() {

	// -P flag to set the port
	// -H flag to set the host
	// -C flag to set the config file path
	// -D or --debug flag to enable debug mode
	isDebug := false
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-P":
			connPort = os.Args[i+1]
			i++
		case "-H":
			connHost = os.Args[i+1]
			i++
		case "-C":
			configPath = os.Args[i+1]
			i++
		case "-D", "--debug":
			isDebug = true
		}
	}

	server.NewServer(connHost, connPort, configPath, isDebug)
}
