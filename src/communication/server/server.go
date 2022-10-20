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

func (server *Server) IsDebug() bool {
	return server.isDebug
}

type Debuggable interface {
	IsDebug() bool
}

func Debug(source Debuggable, message string) {
	if source.IsDebug() {
		fmt.Println("DEBUG: ", message)
	}
}

type Server struct {
	clients           map[int]*Client
	dbm               *DatabaseManager
	messagingProtocol protocol.TcpProtocol
	isDebug           bool
}

func NewServer(isDebug bool) *Server {
	path, err := filepath.Abs(CONFIG_FILE_PATH)

	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}

	return &Server{
		clients:           make(map[int]*Client),
		dbm:               NewDatabaseManager(models.LoadDatabaseFromJson(path), isDebug),
		messagingProtocol: protocol.TcpProtocol{},
		isDebug:           isDebug,
	}
}

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
		newClient := NewClient(server, &conn)
		// Handle connections in a new goroutine.
		go server.handleRequest(newClient)
	}
}

func (server *Server) getNextClientId() uint {
	return uint(len(server.clients))
}

func (server *Server) closeRequest(client *Client) {
	Debug(server, "Closing client connection")
	client.Close()
}

// Handles incoming requests.
func (server *Server) handleRequest(client *Client) {
	Debug(server, "Now we dialogue with client")
	defer server.closeRequest(client)

	for {
		data, err := client.Read()
		if err != nil {
			if err == io.EOF {
				Debug(server, "Client disconnected")
				break
			} else {
				Debug(server, "Error reading from client: "+err.Error())
				break
			}
		}

		if data.Type == protocol.DEBUG && server.isDebug {
			client.isDebug = true
		}

		server.dbm.RequestChannel <- *NewDatabaseRequest(client, data)

	}

}
