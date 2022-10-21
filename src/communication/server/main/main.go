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

	isDebug := false
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-P", "--port":
			connPort = os.Args[i+1]
			i++
		case "-H", "--host":
			connHost = os.Args[i+1]
			i++
		case "-C", "--config":
			configPath = os.Args[i+1]
			i++
		case "-D", "--debug":
			isDebug = true
		}
	}

	server.NewServer(connHost, connPort, configPath, isDebug)
}
