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
	fmt.Print("\n")
	for {
		data, err := reader.ReadString(';')
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		fmt.Println("Data :", data)
		// Remove the semicolon
		data = data[:len(data)-1]
		splitMessage := strings.Split(data, ",")
		code := splitMessage[0]

		if code == LOGIN {
			fmt.Println("user wants to login")
			if len(splitMessage) != 3 {
				fmt.Println("Wrong number of arguments")
				break
			}
			name := splitMessage[1]
			password := splitMessage[2]
			fmt.Print("name: ", name)
			fmt.Println(" password: ", password)
			result, err := login(name, password)
			if err != nil {
				fmt.Println("Error logging in: ", err.Error())
				conn.Write([]byte(NOTOK + ";"))
				break
			}
			if result {
				fmt.Println("Login successful")
				conn.Write([]byte(OK + ";"))
			} else {
				fmt.Println("Login failed")
				conn.Write([]byte(NOTOK + ";"))
			}
		} else if code == CREATE_EVENT {
			fmt.Println("user wants to create an event")
		} else if code == STOP {
			fmt.Println("user wants to stop")
			conn.Close()
			return
		} else {
			fmt.Println("wtf is this")
		}
	}
}
