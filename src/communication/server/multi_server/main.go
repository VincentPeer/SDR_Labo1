package main

import (
	"SDR_Labo1/src/communication/server"
	"fmt"
	"os"
	"strconv"
)

var (
	connHost          = "localhost"
	connPort          = "3333"
	configPath        = "./config.json"
	networkConfigPath = "./networkConfig.json"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Error: missing server id")
		os.Exit(1)
	}

	fmt.Println("Starting server with id", os.Args[1])
	idServ, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Error: server id is not a number")
		os.Exit(1)
	}
	server.StartMultiServer(idServ, networkConfigPath, configPath, true)
}
