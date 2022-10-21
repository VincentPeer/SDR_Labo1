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
	CLOSE_EVENT  = "CLOSE_EVENT"
	DEBUG        = "DEBUG"
	EVENT_REG    = "EVENT_REG"
	GET_EVENTS   = "GET_EVENTS"
	GET_JOBS     = "GET_JOBS"
	STOP         = "STOP"
	DELIMITER    = ';'
	SEPARATOR    = ','
)

// SDRProtocol is the protocol used by the server and the client
type SDRProtocol struct {
}

// ToSend Formats a message to be sent
//
// Returns the message to be sent
func (t *SDRProtocol) ToSend(data DataPacket) (string, error) {
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

// Receive Parses a message received
func (t *SDRProtocol) Receive(message string) (DataPacket, error) {
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

// GetDelimiter returns the delimiter used by the protocol
func (t *SDRProtocol) GetDelimiter() byte {
	return DELIMITER
}

// NewError creates a new error packet
func (t *SDRProtocol) NewError(message string) DataPacket {
	return DataPacket{
		Type: NOTOK,
		Data: []string{message},
	}
}

// NewSuccess creates a new success packet
func (t *SDRProtocol) NewSuccess(message string) DataPacket {
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
