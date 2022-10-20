package server

import (
	"SDR_Labo1/src/communication/protocol"
	"bufio"
	"net"
)

// Client represent a client communicating with the server
type Client struct {
	bufin         *bufio.Reader
	bufout        *bufio.Writer
	conn          *net.Conn
	connectedUser string  // name of the user identified by the client
	server        *Server // server the client is connected to
	isDebug       bool
}

// NewClient creates a new client
func NewClient(server *Server, conn *net.Conn) *Client {
	return &Client{
		bufin:   bufio.NewReader(*conn),
		bufout:  bufio.NewWriter(*conn),
		conn:    conn,
		server:  server,
		isDebug: false,
	}
}

// Read reads a message from the client
func (c *Client) Read() (protocol.DataPacket, error) {
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
func (c *Client) Write(data protocol.DataPacket) error {
	message, err := c.server.messagingProtocol.ToSend(data)

	if err == nil {
		_, err = c.bufout.WriteString(message)
		c.bufout.Flush()
	}
	return err
}

func (c *Client) SendError(message string) error {
	return c.Write(c.server.messagingProtocol.NewError(message))
}

func (c *Client) SendSuccess(message string) error {
	return c.Write(c.server.messagingProtocol.NewSuccess(message))
}

func (c *Client) Close() { // TODO goroutine safe
	(*c.conn).Close()
}

func (c *Client) isLogged() bool {
	return c.connectedUser != ""
}

func (c *Client) Login(username string) {
	c.connectedUser = username
}

func (c *Client) Logout() {
	c.connectedUser = ""
}

func (c *Client) GetConnected() string {
	return c.connectedUser
}
