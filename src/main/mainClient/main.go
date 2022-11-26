package main

import (
	"SDR_Labo1/src/communication/client"
	"os"
)

// Default connection values
var (
	connHost = "localhost"
	connPort = []string{"11211", "11212", "11213"}
)

// Main function that start a new client
func main() {

	isDebug := false
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-P", "--port":
			connPort[0] = os.Args[i+1] // TODO change this to support multiple servers
			i++
		case "-H", "--host":
			connHost = os.Args[i+1]
			i++
		case "-D", "--debug":
			isDebug = true
		}
	}

	conn := client.CreateConnection(connHost, connPort[1], isDebug) // Create a new connection

	client.StartUI(conn) // Start the user interface
	conn.Close()
}
