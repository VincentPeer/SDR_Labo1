package main

import (
	"bufio"
	"errors"
	"net"
	"strings"
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

type client struct {
	//state clientState
	ID     int
	bufin  *bufio.Reader
	bufout *bufio.Writer
	conn   *net.Conn
}

// NewClient creates a new client
func NewClient(id int, conn *net.Conn) *client {
	return &client{
		ID:     id,
		bufin:  bufio.NewReader(*conn),
		bufout: bufio.NewWriter(*conn),
		conn:   conn,
	}
}

// Read reads a message from the client
func (c *client) Read() (string, error) {
	return c.bufin.ReadString(DELIMITER)
}

// Write writes a message to the client
func (c *client) Write(message string) error {
	// Delimiters can't be part of the message
	if strings.Contains(message, string(DELIMITER)) {
		return errors.New("message contains delimiter")
	}

	_, err := c.bufout.WriteString(message + string(DELIMITER))
	c.bufout.Flush()
	return err
}
