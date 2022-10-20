package main

import (
	"SDR_Labo1/src/communication/client"
	"os"
)

// Main function that start a new client
func main() {

	// If the program is started with -d as the first argument, the client will run in debug mode
	debug := false
	if len(os.Args) > 1 && os.Args[1] == "-d" {
		debug = true
	}

	conn := client.CreateConnection(debug) // Establish a new connection with the server

	client.StartUI(conn) // Start the user interface

}
