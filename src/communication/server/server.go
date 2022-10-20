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

type Server struct {
	clients           map[int]*Client
	dbm               *DatabaseManager
	messagingProtocol protocol.TcpProtocol
}

func NewServer() *Server {
	path, err := filepath.Abs(CONFIG_FILE_PATH)

	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}

	return &Server{
		clients:           make(map[int]*Client),
		dbm:               NewDatabaseManager(models.LoadDatabaseFromJson(path)),
		messagingProtocol: protocol.TcpProtocol{},
	}
}

func (server *Server) Start() {

	go server.dbm.Start()
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
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

func closeRequest(client *Client) {
	fmt.Println("Closing client")
	client.Close()
}

// Handles incoming requests.
func (server *Server) handleRequest(client *Client) {
	fmt.Println("Now we dialogue with client")
	defer closeRequest(client)

	for {
		data, err := client.Read()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				break
			} else {
				fmt.Println("Error reading:", err.Error())
				break
			}
		}

		fmt.Println("Data :", data)
		server.dbm.RequestChannel <- *NewDatabaseRequest(client, data)

	}

}
