package main

import "SDR_Labo1/src/communication/client"

// Main function that start a new client
func main() {
	conn := client.CreateConnection() // Establish a new connection with the server
	client.UserInterface(conn)        // Start the user interface
}
