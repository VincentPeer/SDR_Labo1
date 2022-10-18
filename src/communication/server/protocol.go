package main

type dataPacket struct {
	// The type of the packet
	Type string
	// The data of the packet
	Data []string
}

type protocol interface {
	// Formats a message to be sent to the client
	// Returns the message to be sent
	ToSend(data dataPacket) (string, error)

	// Parses a message received from the client
	// Returns the data to be processed
	Receive(message string) (dataPacket, error)

	// Get the delimiter used by the protocol
	GetDelimiter() byte
}
