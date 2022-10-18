package main

import (
	"errors"
	"strings"
)

const (
	// Protocol
	OK           = "OK"
	NOTOK        = "NOTOK"
	LOGIN        = "LOGIN"
	CREATE_EVENT = "CREATE_EVENT"
	STOP         = "STOP"
	DELIMITER    = ';'
	SEPARATOR    = ','
)

type tcpProtocol struct {
}

// Formats a message to be sent to the client
// Returns the message to be sent
func (t *tcpProtocol) ToSend(data dataPacket) (string, error) {
	fullpacket := append([]string{data.Type}, data.Data...)
	for s := range fullpacket {
		if strings.ContainsRune(fullpacket[s], DELIMITER|SEPARATOR) {
			return "", errors.New("delimiter or separator found in packet")
		}
	}
	if len(data.Data) > 0 {
		return strings.Join(fullpacket, string(SEPARATOR)) + string(DELIMITER), nil
	}
	return data.Type + string(DELIMITER), nil
}

// Parses a message received from the client
// Returns the data to be processed
func (t *tcpProtocol) Receive(message string) (dataPacket, error) {
	// Remove the delimiter
	message = message[:len(message)-1]

	// Split the message into type and data
	split := strings.Split(message, string(SEPARATOR))

	// Return the data
	return dataPacket{
		Type: split[0],
		Data: split[1:],
	}, nil
}

// Get the delimiter used by the protocol
func (t *tcpProtocol) GetDelimiter() byte {
	return DELIMITER
}

func (t *tcpProtocol) NewError(message string) dataPacket {
	return dataPacket{
		Type: NOTOK,
		Data: []string{message},
	}
}

func (t *tcpProtocol) NewSuccess(message string) dataPacket {
	if message == "" {
		return dataPacket{
			Type: OK,
			Data: []string{},
		}
	} else {
		return dataPacket{
			Type: OK,
			Data: []string{message},
		}
	}
}
