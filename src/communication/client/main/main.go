package main

import (
	"SDR_Labo1/src/communication/client"
	"os"
)

func main() {

	// If the program is started with -d as the first argument, the client will run in debug mode
	debug := false
	if len(os.Args) > 1 && os.Args[1] == "-d" {
		debug = true
	}

	conn := client.CreateConnection(debug)
	client.UserInterface(conn)

	//conn.Close()
}
