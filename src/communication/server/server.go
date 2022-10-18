package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

var (
	nbClients         int = 0
	users             []User
	events            []Event
	messagingProtocol = &tcpProtocol{}
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
		nbClients++
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(NewClient(nbClients, &conn, messagingProtocol))
	}
}

// Handles incoming requests.
func handleRequest(client *client) {
	fmt.Println("Now we dialogue with client")

	for {
		data, err := client.Read()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				break
			} else {
				fmt.Println("Error reading:", err.Error())
				continue
			}
		}

		fmt.Println("Data :", data)

		switch data.Type {
		case LOGIN:
			if len(data.Data) != 2 {
				client.Write(messagingProtocol.NewError("Invalid number of arguments"))
				continue
			}

			name := data.Data[0]
			password := data.Data[1]

			fmt.Println("user wants to login")
			fmt.Print("name: ", name)
			fmt.Println(" password: ", password)

			result, err := login(name, password)
			if err != nil {
				fmt.Println("Error logging in: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			}
			if result {
				fmt.Println("Login successful")
				client.Write(messagingProtocol.NewSuccess(""))
			} else {
				fmt.Println("Login failed")
				client.Write(messagingProtocol.NewError("Login failed"))
			}
		case CREATE_EVENT:
			fmt.Println("user wants to create an event")

			if len(data.Data) < 3 {
				client.Write(messagingProtocol.NewError("Invalid number of arguments"))
				continue
			}

			eventName := data.Data[0]
			organizerName := data.Data[1]
			password := data.Data[2]

			result, err := login(organizerName, password)
			if err != nil {
				fmt.Println("Error logging in: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			}
			if !result {
				fmt.Println("Login failed")
				client.Write(messagingProtocol.NewError("Login failed"))
				break
			}

			createEvent(events, eventName, organizerName)
		case STOP:
			fmt.Println("user wants to stop the server")
			client.Close()
		default:
			fmt.Println("Unknown command")
		}
	}
}
