package main

import (
	"SDR_Labo1/src/communication/client"
	"SDR_Labo1/src/communication/server"
	"fmt"
)

func main() {
	RunTest()
}

func RunTest() {
	fmt.Println("Running integration tests")

	client.Createclient()
	server.CreateServer()
}

func loginClient() {
	fmt.Println("Checking login")
}
