// name of the package
package client

import (
	"SDR_Labo1/src/communication/protocol"
	"bufio"
	"fmt"
	"net"
	"os"
)

/* constants needed to connect to the server */
const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
	OK        = "OK"
)

// connection contains buffered readers and writers for the console and the server
type connection struct {
	consoleIn *bufio.Reader
	serverIn  *bufio.Reader
	serverOut *bufio.Writer
	protocol  protocol.Protocol
}

// newConnection creates a new connection with our server
func newConnection(consoleIn *bufio.Reader, serverIn *bufio.Reader, serverOut *bufio.Writer, protocol protocol.Protocol) *connection {

	return &connection{
		consoleIn: consoleIn,
		serverIn:  serverIn,
		serverOut: serverOut,
		protocol:  protocol,
	}

}

// Prepare the connection and start a client
func Createclient() {

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

	// Start the client
	userInterface(client)
}
