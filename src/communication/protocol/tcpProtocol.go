package protocol

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
	JOIN_EVENT   = "JOIN_EVENT"
	GET_EVENTS   = "GET_EVENTS"
	STOP         = "STOP"
	DELIMITER    = ';'
	SEPARATOR    = ','
)

type TcpProtocol struct {
}

// Formats a message to be sent to the client
// Returns the message to be sent
func (t *TcpProtocol) ToSend(data DataPacket) (string, error) {
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
func (t *TcpProtocol) Receive(message string) (DataPacket, error) {
	// Remove the delimiter
	message = message[:len(message)-1]

	// Split the message into type and data
	split := strings.Split(message, string(SEPARATOR))

	// Return the data
	return DataPacket{
		Type: split[0],
		Data: split[1:],
	}, nil
}

// Get the delimiter used by the protocol
func (t *TcpProtocol) GetDelimiter() byte {
	return DELIMITER
}

func (t *TcpProtocol) NewError(message string) DataPacket {
	return DataPacket{
		Type: NOTOK,
		Data: []string{message},
	}
}

func (t *TcpProtocol) NewSuccess(message string) DataPacket {
	if message == "" {
		return DataPacket{
			Type: OK,
			Data: []string{},
		}
	} else {
		return DataPacket{
			Type: OK,
			Data: []string{message},
		}
	}
}
