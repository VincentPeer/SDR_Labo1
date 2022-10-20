// This package set up the connection on the client side
package client

import (
	"SDR_Labo1/src/communication/protocol"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

// constants needed to connect to the a
const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
	OK        = "OK"
)

// connection contains buffered readers and writers to allow communication between the client and the a
type connection struct {
	consoleIn *bufio.Reader
	serverIn  *bufio.Reader
	serverOut *bufio.Writer
	protocol  protocol.Protocol
}

// newConnection establishes a new connection based on our own protocol
func NewConnection(consoleIn *bufio.Reader, serverIn *bufio.Reader, serverOut *bufio.Writer, protocol protocol.Protocol) *connection {
	return &connection{
		consoleIn: consoleIn,
		serverIn:  serverIn,
		serverOut: serverOut,
		protocol:  protocol,
	}
}

// Prepare the connection and start a client
func Createclient() *connection {

	conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	consoleReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)
	serverWriter := bufio.NewWriter(conn)
	client := NewConnection(consoleReader, serverReader, serverWriter, &protocol.TcpProtocol{})

	// Start the client
	userInterface(client)

	return client
}

// LoginClient ask the user to enter his username and password and check if the login is correct
// the client is asked to enter his username and password until the login is correct
func (c *connection) LoginClient(username string, password string) bool {
	// Send the login request to the a
	login := protocol.DataPacket{Type: protocol.LOGIN, Data: []string{username, password}}

	// Manage the login request with the a
	response, _ := c.serverRequest(login)
	return response
}

// createEvent creates a new event makde by an organizer
func (c *connection) createEvent(jobList []string) bool {
	event := protocol.DataPacket{Type: protocol.CREATE_EVENT, Data: jobList}
	response, _ := c.serverRequest(event)
	return response
}

func (c *connection) printEvents() {
	eventFound, data := c.serverRequest(protocol.DataPacket{Type: protocol.GET_EVENTS})

	if eventFound {
		printDataPacket(data)
	}
}

func (c *connection) volunteerRegistration(eventId int, jobId int) {
	request := protocol.DataPacket{Type: protocol.EVENT_REG, Data: []string{strconv.Itoa(eventId), strconv.Itoa(jobId)}}
	c.serverRequest(request)
}

func (c *connection) listJobs(eventId int) {
	request := protocol.DataPacket{Type: protocol.GET_EVENTS, Data: []string{strconv.Itoa(eventId)}}
	response, data := c.serverRequest(request)

	if response {
		printDataPacket(data)
	}
}

func (c *connection) volunteerRepartition(eventId int) {
	eventId = c.integerReader("Enter event id : ")
	request := protocol.DataPacket{Type: protocol.GET_JOBS, Data: []string{strconv.Itoa(eventId)}}
	response, data := c.serverRequest(request)

	if response {
		printDataPacket(data)
	}
}

func (c *connection) closeEvent(eventId int) {
	closeEvent := protocol.DataPacket{Type: protocol.CLOSE_EVENT, Data: []string{strconv.Itoa(eventId)}}
	c.serverRequest(closeEvent)
}

func (c *connection) readFromServer() protocol.DataPacket {
	message, err := c.serverIn.ReadString(protocol.DELIMITER)
	if err != nil {
		log.Fatal(err)
	}
	data, e := messagingProtocol.Receive(message)
	if e != nil {
		log.Fatal(e)
	}
	return data
}

func (c *connection) writeToServer(data protocol.DataPacket) {
	m, e := messagingProtocol.ToSend(data)
	if e != nil {
		log.Fatal(e)
	}
	_, writtingError := c.serverOut.WriteString(m)
	if writtingError != nil {
		log.Fatal(writtingError)
	}
	flushError := c.serverOut.Flush()
	if flushError != nil {
		log.Fatal(flushError)
	}
}

func (c *connection) serverRequest(data protocol.DataPacket) (bool, protocol.DataPacket) {
	c.writeToServer(data)
	response := c.readFromServer()
	if response.Type != protocol.OK {
		fmt.Println(response.Data)
	}
	return response.Type == protocol.OK, response
}
