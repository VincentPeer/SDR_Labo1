package main

import (
	"SDR_Labo1/src/ui"
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {

	connection, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Connection done!")
	}
	defer connection.Close()

	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(connection)
	writer := bufio.NewWriter(connection)

	ui.UserInterface(reader, writer, serverReader)

}
