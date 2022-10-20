package main

import (
	"SDR_Labo1/src/communication/server"
	"os"
)

func main() {
	// If the program is started with -d as the first argument, the server will run in debug mode
	debug := false
	if len(os.Args) > 1 && os.Args[1] == "-d" {
		debug = true
	}

	server := server.NewServer(debug)
	server.Start()
}
