/*
This package sets up the TCP server and handles incoming requests.
*/
package server

import (
	"SDR_Labo1/src/communication/protocol"
	"SDR_Labo1/src/communication/server/models"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

const (
	CONN_HOST        = "localhost"
	CONN_PORT        = "3333"
	CONN_TYPE        = "tcp"
	CONFIG_FILE_PATH = "./config.json"
)

// Server listens for incoming connections on a given port and forwards them to the database manager
type Server struct {
	dbm               *DatabaseManager
	messagingProtocol protocol.SDRProtocol
	isDebug           bool
}

// NewServer returns a ready to use TCP server
//
// If debug is true, the server will print debug messages
func NewServer(isDebug bool) *Server {
	path, err := filepath.Abs(CONFIG_FILE_PATH)

	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}

	return &Server{
		dbm:               NewDatabaseManager(models.LoadDatabaseFromJson(path), isDebug),
		messagingProtocol: protocol.SDRProtocol{},
		isDebug:           isDebug,
	}
}

// IsDebug returns true if the server is in debug mode
func (server *Server) IsDebug() bool {
	return server.isDebug
}

// Start listening for incoming connections. Will block unless an error occurs.
func (server *Server) Start() {

	go server.dbm.Start()
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		Debug(server, "Error listening: "+err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	Debug(server, "Listening on "+CONN_HOST+":"+CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			Debug(server, "Error accepting: "+err.Error())
			os.Exit(1)
		}
		newClient := newClientConnection(server, &conn)
		// Handle connections in a new goroutine.
		go server.handleRequest(newClient)
	}
}

// closeRequest closes the connection with the client and removes it from the list of clients
func (server *Server) closeRequest(client *clientConnection) {
	Debug(server, "Closing client connection")
	client.close()
}

// handleRequest handles incoming requests from clients and forwards database access requests to the database manager
func (server *Server) handleRequest(client *clientConnection) {
	Debug(server, "Now we dialogue with client")
	defer server.closeRequest(client)

	for {
		data, err := client.read()
		if err != nil {
			if err == io.EOF { // Client disconnected
				Debug(server, "Client disconnected")
				break
			} else {
				Debug(server, "Error reading from client: "+err.Error())
				break
			}
		}

		if data.Type == protocol.DEBUG && server.isDebug { // Client sent a debug command so we set him to debug mode
			client.isDebug = true
		}

		server.dbm.RequestChannel <- *NewDatabaseRequest(client, data) // Forward the request to the database manager

	}

}
