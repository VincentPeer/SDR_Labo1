// This package set up the Connection on the client side
package client

import (
	"SDR_Labo1/src/communication/protocol"
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

// constants needed to connect to the a
const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
	OK        = "OK"
)

// Connection contains buffered readers and writers to allow communication between the client and the a
type Connection struct {
	conn      net.Conn
	serverIn  *bufio.Reader
	serverOut *bufio.Writer
	protocol  protocol.Protocol
}

// newConnection establishes a new Connection based on our own protocol
func NewConnection(conn net.Conn, protocol protocol.Protocol) *Connection {
	return &Connection{
		serverIn:  bufio.NewReader(conn),
		serverOut: bufio.NewWriter(conn),
		protocol:  protocol,
	}
}

// Prepare the Connection and start a client
func CreateConnection() *Connection {

	conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Fatal(err)
	}

	client := NewConnection(conn, &protocol.TcpProtocol{})

	return client
}

// LoginClient ask the user to enter his username and password and check if the login is correct
// the client is asked to enter his username and password until the login is correct
func (c *Connection) LoginClient(username string, password string) bool {
	// Send the login request to the server
	login := protocol.DataPacket{Type: protocol.LOGIN, Data: []string{username, password}}

	response, _ := c.serverRequest(login)
	return response
}

// createEvent creates a new event makde by an organizer
func (c *Connection) CreateEvent(jobList []string) bool {
	event := protocol.DataPacket{Type: protocol.CREATE_EVENT, Data: jobList}
	response, _ := c.serverRequest(event)
	return response
}

func (c *Connection) PrintEvents() {
	eventFound, data := c.serverRequest(protocol.DataPacket{Type: protocol.GET_EVENTS})

	if eventFound {
		printDataPacket(data)
	}
}

func (c *Connection) VolunteerRegistration(eventId int, jobId int) {
	request := protocol.DataPacket{Type: protocol.EVENT_REG, Data: []string{strconv.Itoa(eventId), strconv.Itoa(jobId)}}
	c.serverRequest(request)
}

func (c *Connection) ListJobs(eventId int) {
	request := protocol.DataPacket{Type: protocol.GET_EVENTS, Data: []string{strconv.Itoa(eventId)}}
	response, data := c.serverRequest(request)

	if response {
		printDataPacket(data)
	}
}

func (c *Connection) VolunteerRepartition(eventId int) {
	request := protocol.DataPacket{Type: protocol.GET_JOBS, Data: []string{strconv.Itoa(eventId)}}
	response, data := c.serverRequest(request)

	if response {
		printDataPacket(data)
	}
}

func (c *Connection) CloseEvent(eventId int) {
	closeEvent := protocol.DataPacket{Type: protocol.CLOSE_EVENT, Data: []string{strconv.Itoa(eventId)}}
	c.serverRequest(closeEvent)
}

func (c *Connection) readFromServer() protocol.DataPacket {
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

func (c *Connection) writeToServer(data protocol.DataPacket) {
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

func (c *Connection) serverRequest(data protocol.DataPacket) (bool, protocol.DataPacket) {
	c.writeToServer(data)
	response := c.readFromServer()
	if response.Type != protocol.OK {
		fmt.Println(response.Data)
	}
	return response.Type == protocol.OK, response
}

func (c *Connection) Close() {
	c.writeToServer(protocol.DataPacket{Type: protocol.STOP})
	if c != nil && c.conn != nil {
		c.conn.Close()
	}
}
