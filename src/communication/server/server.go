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
	connType = "tcp"
)

// Server listens for incoming connections on a given port and forwards them to the database manager
type Server struct {
	Id                int
	host              string
	port              string
	dbm               *databaseManager
	messagingProtocol protocol.SDRProtocol
	debugFlag         bool
}

// NewServer returns a ready to use TCP server and starts it
//
// If debug is true, the server will print debug messages
func NewServer(id int, host string, port string, configFilePath string, isDebug bool, startDirectly bool) *Server {
	path, err := filepath.Abs(configFilePath)

	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}

	srv := &Server{
		Id:                id,
		host:              host,
		port:              port,
		dbm:               newDatabaseManager(models.LoadDatabaseFromJson(path), isDebug),
		messagingProtocol: protocol.SDRProtocol{},
		debugFlag:         isDebug,
	}
	if startDirectly {
		srv.start()
	}
	return srv
}

// IsDebug returns true if the server is in debug mode
func (server *Server) isDebug() bool {
	return server.debugFlag
}

// start listening for incoming connections. Will block unless an error occurs.
func (server *Server) start() {

	go server.dbm.start()
	// Listen for incoming connections.
	l, err := net.Listen(connType, server.host+":"+server.port)
	if err != nil {
		debug(server, "Error listening: "+err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	debug(server, "Listening on "+server.host+":"+server.port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			debug(server, "Error accepting: "+err.Error())
			os.Exit(1)
		}
		newClient := newClientConnection(server, &conn)
		// Handle connections in a new goroutine.
		go server.handleRequest(newClient)
	}
}

// closeRequest closes the connection with the client and removes it from the list of clients
func (server *Server) closeRequest(client *clientConnection) {
	debug(server, "Closing client connection")
	client.close()
}

// handleRequest handles incoming requests from clients and forwards database access requests to the database manager
func (server *Server) handleRequest(client *clientConnection) {
	debug(server, "Now we dialogue with client")
	defer server.closeRequest(client)

	for {
		data, err := client.read()
		if err != nil {
			if err == io.EOF { // Client disconnected
				debug(server, "Client disconnected")
				break
			} else {
				debug(server, "Error reading from client: "+err.Error())
				break
			}
		}

		if data.Type == protocol.DEBUG && server.debugFlag { // Client sent a debug command so we set him to debug mode
			client.isDebug = true
		}

		server.dbm.requestChannel <- *newDatabaseRequest(client, data) // Forward the request to the database manager

	}

}
