package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	CONN_HOST        = "localhost"
	CONN_PORT        = "3333"
	CONN_TYPE        = "tcp"
	CONFIG_FILE_PATH = "./config.json"
)

var (
	users  []User
	events []Event
)

func loadConfig(jsonPath string) Config {
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully opened " + jsonFile.Name())
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var conf Config
	json.Unmarshal(byteValue, &conf)

	return conf
}

func main() {
	path, err := filepath.Abs(CONFIG_FILE_PATH)

	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}

	config := loadConfig(path)
	users = config.Users
	events = config.Events

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
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	reader := bufio.NewReader(conn)
	//writer := bufio.NewWriter(conn)
	fmt.Println("Now we dialogue with client")

	for {
		message, _ := reader.ReadString('\n')
		fmt.Print("user is " + message)
		if strings.Compare(message, "STOP") == 0 {
			fmt.Println("TCP client exiting...")
			return
		}
	}

	/*	// Make a buffer to hold incoming data.
		buf := make([]byte, 1024)
		// Read the incoming connection into the buffer.
		reqLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		if reqLen > 0 {
			fmt.Println("Received data:", string(buf))
		}
		// Send a response back to person contacting us.
		conn.Write([]byte("Message received."))
		// Close the connection when you're done with it.
		conn.Close()*/
}
