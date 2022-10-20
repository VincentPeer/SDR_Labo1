package main

import (
	"SDR_Labo1/src/communication/client"
	"SDR_Labo1/src/communication/server"
	"fmt"
)

func main() {
	fmt.Println("Running integration tests")

	go server.CreateServer()
	conn := client.CreateConnection()
	fmt.Println("Client created...")
	fmt.Println("Starting tests...")
	loginStat := conn.LoginClient("a", "1")
	fmt.Println(loginStat)
	loginClient()

	conn.Close()
}

func loginClient() {
	fmt.Println("Checking login")

}
