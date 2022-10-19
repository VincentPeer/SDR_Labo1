package main

import (
	"SDR_Labo1/src/communication/protocol"
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
	OK        = "OK"
)

type connection struct {
	consoleIn *bufio.Reader
	serverIn  *bufio.Reader
	serverOut *bufio.Writer
	protocol  protocol.Protocol
}

func newConnection(consoleIn *bufio.Reader, serverIn *bufio.Reader, serverOut *bufio.Writer, protocol protocol.Protocol) *connection {

	return &connection{
		consoleIn: consoleIn,
		serverIn:  serverIn,
		serverOut: serverOut,
		protocol:  protocol,
	}

}

func main() {

	conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Connection done!")
	}
	defer conn.Close()

	consoleReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)
	serverWriter := bufio.NewWriter(conn)
	client := newConnection(consoleReader, serverReader, serverWriter, &protocol.TcpProtocol{})

	userInterface(client)

}
