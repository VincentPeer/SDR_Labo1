package server

import (
	"SDR_Labo1/src/communication/protocol"
	"bufio"
	"net"
)

// clientConnection represent a client communicating with the server
type clientConnection struct {
	bufin         *bufio.Reader
	bufout        *bufio.Writer
	conn          *net.Conn
	connectedUser string  // name of the user identified by the client
	server        *Server // server the client is connected to
	isDebug       bool
}

// newClientConnection creates a new client connection
func newClientConnection(server *Server, conn *net.Conn) *clientConnection {
	return &clientConnection{
		bufin:   bufio.NewReader(*conn),
		bufout:  bufio.NewWriter(*conn),
		conn:    conn,
		server:  server,
		isDebug: false,
	}
}

// read reads data the client sent and parse it to a DataPacket object.
func (c *clientConnection) read() (protocol.DataPacket, error) {
	message, err := c.bufin.ReadString(c.server.messagingProtocol.GetDelimiter())

	if c.server.isDebug() {
		debug(c.server, "<-- "+message)
	}

	if err != nil {
		return protocol.DataPacket{}, err
	} else {
		return c.server.messagingProtocol.Receive(message)
	}
}

// write writes a message to the client by converting the DataPacket object to a string.
func (c *clientConnection) write(data protocol.DataPacket) error {
	message, err := c.server.messagingProtocol.ToSend(data)

	if err == nil {
		_, err = c.bufout.WriteString(message)
		c.bufout.Flush()
	}

	if c.server.isDebug() {
		debug(c.server, "--> "+message)
	}
	return err
}

// sendError format an error as a packet and send it to the client.
func (c *clientConnection) sendError(message string) error {
	return c.write(c.server.messagingProtocol.NewError(message))
}

// sendSuccess format a success packet and send it to the client.
func (c *clientConnection) sendSuccess(message string) error {
	return c.write(c.server.messagingProtocol.NewSuccess(message))
}

// close closes the client connection properly.
func (c *clientConnection) close() {
	(*c.conn).Close()
}
