package main

import "SDR_Labo1/src/communication/client"

func main() {
	conn := client.CreateConnection()
	client.UserInterface(conn)

	//conn.Close()
}
