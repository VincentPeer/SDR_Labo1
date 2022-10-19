package main

import (
	"SDR_Labo1/src/communication/protocol"
	"SDR_Labo1/src/communication/server/models"
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
	db                database
	messagingProtocol = &protocol.TcpProtocol{}
)

type database struct {
	Users  models.Users  `json:"users"`
	Events models.Events `json:"events"`
}

func loadConfig(jsonPath string) database {
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully opened " + jsonFile.Name())
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var conf database
	json.Unmarshal(byteValue, &conf)

	return conf
}

func main() {
	path, err := filepath.Abs(CONFIG_FILE_PATH)

	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}

	db = loadConfig(path)

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

func closeRequest(client *client) {
	fmt.Println("Closing client")
	nbClients--
	client.Close()
}

// Handles incoming requests.
func handleRequest(client *client) {
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

		switch data.Type {
		case protocol.LOGIN:
			if len(data.Data) != 2 {
				client.Write(messagingProtocol.NewError("Invalid number of arguments"))
				continue
			}

			name := data.Data[0]
			password := data.Data[1]

			fmt.Println("user wants to login")
			fmt.Print("name: ", name)
			fmt.Println(" password: ", password)

			result, err := db.Users.Login(name, password)
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
		case protocol.CREATE_EVENT:
			fmt.Println("user wants to create an event")

			if len(data.Data) < 3 {
				client.Write(messagingProtocol.NewError("Invalid number of arguments"))
				continue
			}

			eventName := data.Data[0]
			organizerName := data.Data[1]
			password := data.Data[2]

			result, err := db.Users.Login(organizerName, password)
			if err != nil {
				fmt.Println("Error logging in: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			}
			if !result {
				fmt.Println("Login failed")
				client.Write(messagingProtocol.NewError("Login failed"))
				continue
			}

			db.Events.CreateEvent(eventName, organizerName)
		case protocol.STOP:
			fmt.Println("user wants to stop the server")
		default:
			fmt.Println("Unknown command")
		}
	}
}
