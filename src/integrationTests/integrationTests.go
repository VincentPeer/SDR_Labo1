package main

import (
	"SDR_Labo1/src/communication/client"
	"SDR_Labo1/src/communication/server"
	"fmt"
)

func main() {
	fmt.Println("Running integration tests")

	server.CreateServer()
	conn := client.Createclient()
	conn.LoginClient("admin", "admin")
	loginClient()
}

func loginClient() {
	fmt.Println("Checking login")

}
