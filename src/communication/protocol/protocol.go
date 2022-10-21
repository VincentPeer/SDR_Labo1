/*
This package provides the protocol used to communicate between the client and the server
*/
package protocol

// DataPacket is a packet of data sent between the server and the client
type DataPacket struct {
	// The type of the packet
	Type string
	// The data of the packet
	Data []string
}

// Protocol is the protocol used by the server and the client
//
// It is used to send and receive messages
type Protocol interface {
	// Formats a message to be sent to the client
	// Returns the message to be sent
	ToSend(data DataPacket) (string, error)

	// Parses a message received from the client
	// Returns the data to be processed
	Receive(message string) (DataPacket, error)

	// Get the delimiter used by the protocol
	GetDelimiter() byte
}
