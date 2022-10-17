package main

import "bufio"

const (
	// Protocol
	OK           = "OK"
	NOTOK        = "NOTOK"
	LOGIN        = "LOGIN"
	CREATE_EVENT = "CREATE_EVENT"
	STOP         = "STOP"
	DELIMITER    = ';'
)

func sendError(destination *bufio.Writer, message string) error {
	if message != "" {
		message = "," + message
	}
	_, err := destination.WriteString(NOTOK + message + string(DELIMITER))
	destination.Flush()
	return err
}

func sendAck(destination *bufio.Writer, message string) error {
	if message != "" {
		message = "," + message
	}
	_, err := destination.WriteString(OK + message + string(DELIMITER))
	destination.Flush()
	return err
}
