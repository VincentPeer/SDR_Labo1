package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

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

func (s *Server) connectToServer(serverConfig serverConfig) {
	fmt.Println("Connecting to server", serverConfig.Id, "at", serverConfig.Host+":"+serverConfig.Port)
	conn, err := net.Dial(connType, serverConfig.Host+":"+serverConfig.Port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to " + conn.RemoteAddr().String())
	defer conn.Close()
	for {

	}
}

func (s *Server) handleConnectionFromServer(conn *net.Conn) {
	for {

	}
}

func (s *Server) connectToPrecedingServers(config networkConfig) {
	for i := s.Id - 1; i >= 0; i-- {
		go s.connectToServer(config.Servers[i])
	}
}

func (s *Server) startWaitingForServers() {
	l, err := net.Listen(connType, s.host+":"+s.port)
	if err != nil {
		debug(s, "Error listening: "+err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + s.host + ":" + s.port)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			debug(s, "Error accepting: "+err.Error())
			os.Exit(1)
		}
		fmt.Println(s, "New connection from "+conn.RemoteAddr().String())
		go s.handleConnectionFromServer(&conn)
	}
}

func StartMultiServer(id int, serverConfigPath string, dataConfigPath string, isDebug bool) {
	config := ReadNetworkConfig(serverConfigPath)
	if id >= config.NbServers || id < 0 {
		fmt.Println("Error: server id is out of range")
		os.Exit(1)
	}
	server := NewServer(id, config.Servers[id].Host, config.Servers[id].Port, dataConfigPath, isDebug, false)
	go server.startWaitingForServers()
	server.connectToPrecedingServers(config)
	for {

	}
}
