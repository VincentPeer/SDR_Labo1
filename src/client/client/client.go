package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	fmt.Println("Welcome! ")

	fmt.Println("Login : \n Please enter your username : ")
	inputReader := bufio.NewReader(os.Stdin)
	username, _ := inputReader.ReadString('\n')
	fmt.Println("Login : \n Please enter your password : ")
	password, _ := inputReader.ReadString('\n')

	/*	arguments := os.Args
		if len(arguments) == 1 {
			fmt.Println("Please provide host:port.")
			return
		}*/
	//CONNECT := arguments[1]

	connection, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = connection.Write([]byte((username) + "," + password))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(connection, text+"\n")

		message, _ := bufio.NewReader(connection).ReadString('\n')
		fmt.Print("->: " + message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}
