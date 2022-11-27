/*
Package client set up the Connection on the client side and provides methods to communicate with the server
*/
package client

import (
	"SDR_Labo1/src/communication/protocol"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
)

const connType = "tcp"

// Connection contains buffered readers and writers to allow communication between the client and the server
type Connection struct {
	conn      net.Conn
	serverIn  *bufio.Reader
	serverOut *bufio.Writer
	protocol  protocol.Protocol
}

// NewConnection establishes a new Connection based on our own protocol
func NewConnection(conn net.Conn, protocol protocol.Protocol) *Connection {
	return &Connection{
		serverIn:  bufio.NewReader(conn),
		serverOut: bufio.NewWriter(conn),
		protocol:  protocol,
	}
}

// CreateConnection prepare the Connection and start a client
func CreateConnection(host string, port string, isDebug bool) *Connection {
	conn, err := net.Dial(connType, host+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	client := NewConnection(conn, &protocol.SDRProtocol{})
	if isDebug {
		client.sendDebugRequest()
	}

	return client
}

// StartUI allows the user to interact with the server
func StartUI(c *Connection) {
	userInterface(c)
}

// sendDebugRequest sends a debug request to the server
func (c *Connection) sendDebugRequest() {
	c.writeToServer(protocol.DataPacket{Type: protocol.DEBUG})
}

// LoginClient checks the login with the server
//
// Returns true if the login is correct, false otherwise
func (c *Connection) LoginClient(username string, password string) bool {
	// Send the login request to the server
	login := protocol.DataPacket{Type: protocol.LOGIN, Data: []string{username, password}}

	response, _ := c.ServerRequest(login)
	return response
}

// CreateEvent asks the server to add a new event
func (c *Connection) CreateEvent(jobList []string) bool {
	event := protocol.DataPacket{Type: protocol.CREATE_EVENT, Data: jobList}
	response, _ := c.ServerRequest(event)
	return response
}

// PrintEvents asks the server the data containing all the events
// If events are found, they are printed
func (c *Connection) PrintEvents() bool {
	eventFound, data := c.ServerRequest(protocol.DataPacket{Type: protocol.GET_EVENTS})
	if eventFound {
		var objmap map[string]interface{}
		//var objmap map[string]json.RawMessage
		err := json.Unmarshal([]byte(data.Data[0]), &objmap)
		if err != nil {
			log.Fatal(err)
		}
		printEventTable(objmap)
		return true
	}
	return false
}

// VolunteerRegistration asks the server to add a new volunteer
func (c *Connection) VolunteerRegistration(eventId int, jobId int) bool {
	request := protocol.DataPacket{Type: protocol.EVENT_REG, Data: []string{strconv.Itoa(eventId), strconv.Itoa(jobId)}}
	response, _ := c.ServerRequest(request)
	return response
}

// ListJobs asks the server the data containing all the jobs for a specific event
//
// If jobs are found, they are printed
func (c *Connection) ListJobs(eventId int) bool {
	request := protocol.DataPacket{Type: protocol.GET_EVENTS, Data: []string{strconv.Itoa(eventId)}}
	response, data := c.ServerRequest(request)

	if response {
		var objmap map[string]interface{}
		//var objmap map[string]json.RawMessage
		err := json.Unmarshal([]byte(data.Data[0]), &objmap)
		if err != nil {
			log.Fatal(err)
		}
		printJobTAble(objmap)
		return true
	}
	return false
}

// VolunteerRepartition asks the server the repartition of volunteers for a specific event
//
// If repartition is found, it is printed
func (c *Connection) VolunteerRepartition(eventId int) bool {
	request := protocol.DataPacket{Type: protocol.GET_JOBS, Data: []string{strconv.Itoa(eventId)}}
	response, data := c.ServerRequest(request)

	if response {
		var objmap map[string]interface{}
		//var objmap map[string]json.RawMessage
		err := json.Unmarshal([]byte(data.Data[0]), &objmap)
		if err != nil {
			log.Fatal(err)
		}
		printJobPrepartitionTable(objmap)
		return true
	}
	return false
}

// CloseEvent asks the server to close an event by specifying its id
func (c *Connection) CloseEvent(eventId int) bool {
	closeEvent := protocol.DataPacket{Type: protocol.CLOSE_EVENT, Data: []string{strconv.Itoa(eventId)}}
	response, _ := c.ServerRequest(closeEvent)
	return response
}

// readFromServer reads a response from the server
//
// It extracts the data from the response and returns it
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

// writeToServer sends a DataPacket to the server
func (c *Connection) writeToServer(data protocol.DataPacket) {
	m, e := messagingProtocol.ToSend(data)
	if e != nil {
		log.Fatal(e)
	}
	_, writingError := c.serverOut.WriteString(m)
	if writingError != nil {
		log.Fatal(writingError)
	}
	flushError := c.serverOut.Flush()
	if flushError != nil {
		log.Fatal(flushError)
	}
}

// ServerRequest sends a request to the server
//
// Returns as first parameter true if the request was successful, false otherwise and
// returns as second parameter the data received from the server
//
// If the request was not successful, we print the error message received from the server
func (c *Connection) ServerRequest(data protocol.DataPacket) (bool, protocol.DataPacket) {
	c.writeToServer(data)
	response := c.readFromServer()
	if response.Type != protocol.OK {
		fmt.Println(response.Data)
	}
	return response.Type == protocol.OK, response
}

// Close terminates the connection with the server
func (c *Connection) Close() {
	c.writeToServer(protocol.DataPacket{Type: protocol.STOP})
	if c != nil && c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}
