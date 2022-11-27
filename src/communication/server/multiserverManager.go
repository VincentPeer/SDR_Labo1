package server

import (
	"SDR_Labo1/src/communication/protocol"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

var (
	serverInitWaitGroup sync.WaitGroup    // Wait for all servers to be connected
	serverListeners     []*serverListener // List of all server listeners
)

// serverListener is used to manage the communication between this server and another server
type serverListener struct {
	peerServer          *clientConnection        // Connection to the other server
	lamportRequestChan  chan protocol.DataPacket // Channel used by the Lamport algorithm to send requests to the other server
	lamportResponseChan chan protocol.DataPacket // Channel used to forward responses from the other server to the Lamport algorithm
	peerServerChan      chan protocol.DataPacket // Channel where the data received from the other server is forwarded
	lamportRegister     *protocol.DataPacket     // Last information saved by the Lamport algorithm
}

// newServerListener creates a new server listener from a connection
func newServerListener(server *clientConnection) serverListener {
	return serverListener{server, make(chan protocol.DataPacket), make(chan protocol.DataPacket), make(chan protocol.DataPacket), nil}
}

// serverConfig is used to read the configuration of a server from a JSON file
type serverConfig struct {
	Id   int    `json:"id"`
	Port string `json:"port"`
	Host string `json:"host"`
}

// networkConfig is used to read the configuration of the network from a JSON file
type networkConfig struct {
	NbServers int            `json:"nbServers"`
	Servers   []serverConfig `json:"servers"`
}

// ReadNetworkConfig reads the network configuration from a JSON file and returns the configuration
func ReadNetworkConfig(path string) networkConfig {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := networkConfig{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}
	return config
}

// talkWith is used to start the communication between this server and another server
// it forwards the data received from the other server to the Lamport algorithm
func (s *Server) talkWith(listener *serverListener) {
	// Create a go routine to forward the data received to a channel
	go func() {
		for {
			data, err := listener.peerServer.read()
			if err != nil {
				if err == io.EOF { // Client disconnected
					debug(s, "Client disconnected")
					break
				} else {
					debug(s, "Error reading from client: "+err.Error())
					break
				}
			}
			listener.peerServerChan <- data
		}
	}()
	for {
		select { // Wait for a message from the other server or from the Lamport algorithm
		case data := <-listener.lamportRequestChan:
			listener.peerServer.write(data) // Forward the message to the other server
		case data := <-listener.peerServerChan:
			{
				debug(s, "Received data from peer server: "+data.Data[0])
				lamportReceive(listener, data) // Handle the message received from the other server
			}
		}
	}
}

// connectToServer connects to a server id in a given network configuration and starts the communication
func (s *Server) connectToServer(networkConfig networkConfig, id int) {
	serverConfig := networkConfig.Servers[id]
	conn, err := net.Dial(connType, serverConfig.Host+":"+serverConfig.Port)
	if err != nil {
		log.Fatal(err)
	}

	// Sends a handshake to the server
	clientServer := newClientConnection(s, &conn)
	clientServer.write(protocol.DataPacket{Type: protocol.REQ, Data: []string{strconv.Itoa(s.Id)}})
	defer clientServer.close()
	// TODO use another type that REQ to establish connection

	// Waits for the other server to accept the connection
	for {
		data, err := clientServer.read()
		if err != nil {
			if err == io.EOF { // Client disconnected
				debug(s, "Client disconnected")
				break
			} else {
				debug(s, "Error reading from client: "+err.Error())
				break
			}
		}
		if data.Type == protocol.ACK {
			serverInitWaitGroup.Done()
			break
		}
	}

	debug(s, "Conversation with server: "+strconv.Itoa(id))
	serverInitWaitGroup.Wait() // Wait for all servers to be connected before starting the communication

	listener := newServerListener(clientServer)
	serverListeners = append(serverListeners, &listener)
	s.talkWith(&listener)
}

// handleConnectionFromServer handles the handshake with a server that wants to connect to this server
func (s *Server) handleConnectionFromServer(conn *net.Conn) {
	clientServer := newClientConnection(s, conn)
	defer clientServer.close()
	data, err := clientServer.read()
	if err != nil {
		if err == io.EOF { // Client disconnected
			debug(s, "Client disconnected")
			return
		} else {
			debug(s, "Error reading from client: "+err.Error())
			return
		}
	}

	// TODO use another type that REQ to establish connection
	if data.Type == protocol.REQ {
		clientServer.write(protocol.DataPacket{Type: protocol.ACK, Data: []string{}})
	}

	debug(s, "Conversation with server: "+data.Data[0])
	serverInitWaitGroup.Done() // Wait for all servers to be connected before starting the communication
	serverInitWaitGroup.Wait()

	listener := newServerListener(clientServer)
	serverListeners = append(serverListeners, &listener)
	s.talkWith(&listener)
}

// connectToPrecedingServers connects to all the servers with a lower id than this server
func (s *Server) connectToPrecedingServers(config networkConfig) {
	for i := s.Id - 1; i >= 0; i-- {
		go s.connectToServer(config, i)
	}
}

// startWaitingForServers wait for nbServers - id - 1 servers to connect to this server
// for each server that connects, it starts the communication in a new go routine
func (s *Server) startWaitingForServers(nbServers int) {
	l, err := net.Listen(connType, s.host+":"+s.port)
	if err != nil {
		debug(s, "Error listening: "+err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	for i := s.Id + 1; i < nbServers; i++ {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			debug(s, "Error accepting: "+err.Error())
			os.Exit(1)
		}
		go s.handleConnectionFromServer(&conn)
	}

	debug(s, "Don't accept any more connections")
}

// StartMultiServer starts a server in a multi server configuration
// It connects to all the servers with a lower id than this server
// It waits for nbServers - id - 1 servers to connect to this server
// serverConfigPath is the path to the configuration file of the network
// dataConfigPath is the path to the configuration file of the starting data
// debug is a boolean that indicates if the server should print debug messages
func StartMultiServer(id int, serverConfigPath string, dataConfigPath string, isDebug bool) {
	config := ReadNetworkConfig(serverConfigPath)
	if id >= config.NbServers || id < 0 {
		fmt.Println("Error: server id is out of range")
		os.Exit(1)
	}
	// Creating the local server
	server := NewServer(id, config.Servers[id].Host, config.Servers[id].Port, dataConfigPath, isDebug, false)
	debug(server, "Starting server:"+strconv.Itoa(id))
	serverInitWaitGroup.Add(config.NbServers - 1)
	go server.startWaitingForServers(config.NbServers) // Start waiting for servers to connect to this server
	server.connectToPrecedingServers(config)           // Connect to all the servers with a lower id than this server
	serverInitWaitGroup.Wait()                         // Wait for all servers to be connected before accepting requests from clients
	debug(server, "All servers connected")
	server.start()
}
