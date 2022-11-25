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
	waitGroup sync.WaitGroup
)

type multiServerThread struct {
	server            *Server
	conn              *net.Conn
	connectionCounter chan int
}

type serverConfig struct {
	Id   int    `json:"id"`
	Port string `json:"port"`
	Host string `json:"host"`
}

type networkConfig struct {
	NbServers int            `json:"nbServers"`
	Servers   []serverConfig `json:"servers"`
}

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

func handleRequestFromServer(conn net.Conn, server *Server) {

}

func (s *Server) connectToServer(networkConfig networkConfig, id int) {
	serverConfig := networkConfig.Servers[id]
	conn, err := net.Dial(connType, serverConfig.Host+":"+serverConfig.Port)
	if err != nil {
		log.Fatal(err)
	}

	// Sends a request to establish connection with other server
	clientServer := newClientConnection(s, &conn)
	clientServer.write(protocol.DataPacket{Type: protocol.REQ, Data: []string{strconv.Itoa(s.Id)}})

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
			waitGroup.Done()
			break
		}

	}

	debug(s, "Conversation with server: "+strconv.Itoa(id))
	waitGroup.Wait()
	// TODO: talk to other servers
}

func (s *Server) handleConnectionFromServer(conn *net.Conn) {
	clientServer := newClientConnection(s, conn)
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

	if data.Type == protocol.REQ {
		clientServer.write(protocol.DataPacket{Type: protocol.ACK, Data: []string{}})
	}
	debug(s, "Conversation with server: "+data.Data[0])
	waitGroup.Done()
}

func (s *Server) connectToPrecedingServers(config networkConfig) {
	for i := s.Id - 1; i >= 0; i-- {
		go s.connectToServer(config, i)
	}
}

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

func StartMultiServer(id int, serverConfigPath string, dataConfigPath string, isDebug bool) {
	config := ReadNetworkConfig(serverConfigPath)
	if id >= config.NbServers || id < 0 {
		fmt.Println("Error: server id is out of range")
		os.Exit(1)
	}
	server := NewServer(id, config.Servers[id].Host, config.Servers[id].Port, dataConfigPath, isDebug, false)
	debug(server, "Starting server:"+strconv.Itoa(id))
	waitGroup.Add(config.NbServers - 1)
	go server.startWaitingForServers(config.NbServers)
	server.connectToPrecedingServers(config)
	waitGroup.Wait()
	debug(server, "All servers connected")
	server.start()
}
