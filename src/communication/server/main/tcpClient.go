package main

import (
	"SDR_Labo1/src/communication/protocol"
	"bufio"
	"net"
)

//type clientState int
//
//const (
//	greeting clientState = iota
//	creatingEvent
//	registeringForEvent
//	registeringForJob
//	shutdown
//)

// client represents a client connected to the server
type client struct {
	//state clientState
	ID       int
	bufin    *bufio.Reader
	bufout   *bufio.Writer
	conn     *net.Conn
	protocol protocol.Protocol
}

// NewClient creates a new client
func NewClient(id int, conn *net.Conn, protocol protocol.Protocol) *client {
	return &client{
		ID:       id,
		bufin:    bufio.NewReader(*conn),
		bufout:   bufio.NewWriter(*conn),
		conn:     conn,
		protocol: protocol,
	}
}

// Read reads a message from the client
func (c *client) Read() (protocol.DataPacket, error) {
	// Read the message
	message, err := c.bufin.ReadString(c.protocol.GetDelimiter())

	if err != nil {
		return protocol.DataPacket{}, err
	} else {
		// Parse the message
		return c.protocol.Receive(message)
	}
}

// Write writes a message to the client
func (c *client) Write(data protocol.DataPacket) error {
	message, err := c.protocol.ToSend(data)

	if err == nil {
		_, err = c.bufout.WriteString(message)
		c.bufout.Flush()
	}
	return err
}

func (c *client) Close() { // TODO goroutine safe
	(*c.conn).Close()
}
