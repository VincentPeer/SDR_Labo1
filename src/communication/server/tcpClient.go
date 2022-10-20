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

// Read reads data the client sent and parse it to a DataPacket object.
func (c *Client) Read() (protocol.DataPacket, error) {
	message, err := c.bufin.ReadString(c.server.messagingProtocol.GetDelimiter())

	if err != nil {
		return protocol.DataPacket{}, err
	} else {
		return c.server.messagingProtocol.Receive(message)
	}
}

// Write writes a message to the client by converting the DataPacket object to a string.
func (c *Client) Write(data protocol.DataPacket) error {
	message, err := c.server.messagingProtocol.ToSend(data)

	if err == nil {
		_, err = c.bufout.WriteString(message)
		c.bufout.Flush()
	}
	return err
}

// SendError format an error as a packet and send it to the client.
func (c *Client) SendError(message string) error {
	return c.Write(c.server.messagingProtocol.NewError(message))
}

// SendSuccess format a success packet and send it to the client.
func (c *Client) SendSuccess(message string) error {
	return c.Write(c.server.messagingProtocol.NewSuccess(message))
}

// Close closes the client connection properly.
func (c *Client) Close() {
	(*c.conn).Close()
}

// isLogged returns true if the client is logged in.
func (c *Client) isLogged() bool {
	return c.connectedUser != ""
}

// Login sets a user as connected
func (c *Client) Login(username string) {
	c.connectedUser = username
}

// Logout sets the client as disconnected.
func (c *Client) Logout() {
	c.connectedUser = ""
}

// GetConnected get the username connected to the client
func (c *Client) GetConnected() string {
	return c.connectedUser
}
