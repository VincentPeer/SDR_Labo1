package protocol

type DataPacket struct {
	// The type of the packet
	Type string
	// The data of the packet
	Data []string
}

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
