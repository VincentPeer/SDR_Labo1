package main

import "bufio"

const (
	// Protocol
	OK           = "OK"
	NOTOK        = "NOTOK"
	LOGIN        = "LOGIN"
	CREATE_EVENT = "CREATE_EVENT"
	STOP         = "STOP"
)

func sendError(destination *bufio.Writer, message string) error {
	_, err := destination.WriteString(NOTOK + " " + message + "\n")
	destination.Flush()
	return err
}

func sendAck(destination *bufio.Writer, message string) error {
	_, err := destination.WriteString(OK + " " + message + "\n")
	destination.Flush()
	return err
}
