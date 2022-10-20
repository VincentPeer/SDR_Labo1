package server

import (
	"SDR_Labo1/src/communication/protocol"
	"bufio"
	"net"
)

type clientState int

const (
	greeting clientState = iota
	connected
	creatingEvent
	registeringForEvent
	registeringForJob
	shutdown
)

// client represents a client connected to the main
type client struct {
	state         clientState
	ID            uint
	bufin         *bufio.Reader
	bufout        *bufio.Writer
	conn          *net.Conn
	connectedUser string
	server        *Server
}

// NewClient creates a new client
func NewClient(server *Server, conn *net.Conn) *client {
	return &client{
		ID:     server.getNextClientId(),
		state:  greeting,
		bufin:  bufio.NewReader(*conn),
		bufout: bufio.NewWriter(*conn),
		conn:   conn,
		server: server,
	}
}

// Read reads a message from the client
func (c *client) Read() (protocol.DataPacket, error) {
	// Read the message
	message, err := c.bufin.ReadString(c.server.messagingProtocol.GetDelimiter())

	if err != nil {
		return protocol.DataPacket{}, err
	} else {
		// Parse the message
		return c.server.messagingProtocol.Receive(message)
	}
}

// Write writes a message to the client
func (c *client) Write(data protocol.DataPacket) error {
	message, err := c.server.messagingProtocol.ToSend(data)

	if err == nil {
		_, err = c.bufout.WriteString(message)
		c.bufout.Flush()
	}
	return err
}

func (c *client) SendError(message string) error {
	return c.Write(c.server.messagingProtocol.NewError(message))
}

func (c *client) SendSuccess(message string) error {
	return c.Write(c.server.messagingProtocol.NewSuccess(message))
}

func (c *client) Close() { // TODO goroutine safe
	(*c.conn).Close()
}

func (c *client) GetState() clientState {
	return c.state
}

func (c *client) Login(username string) {
	c.connectedUser = username
	c.state = connected
}

func (c *client) Logout() {
	c.connectedUser = ""
	c.state = greeting
}

func (c *client) GetConnected() string {
	return c.connectedUser
}
