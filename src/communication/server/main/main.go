package main

import (
	"SDR_Labo1/src/communication/server"
)

func main() {
	server := server.NewServer()
	server.Start()
}
