package main

import (
	"SDR_Labo1/src/communication/protocol"
	"SDR_Labo1/src/communication/server/models"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

const (
	CONN_HOST        = "localhost"
	CONN_PORT        = "3333"
	CONN_TYPE        = "tcp"
	CONFIG_FILE_PATH = "./config.json"
)

var (
	nbClients         int = 0
	db                models.Database
	messagingProtocol = &protocol.TcpProtocol{}
)

func CreateServer() {
	main()
}

func main() {
	path, err := filepath.Abs(CONFIG_FILE_PATH)

	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}

	db = models.LoadDatabaseFromJson(path)

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

			result, err := db.Login(name, password)
			if err != nil {
				fmt.Println("Error logging in: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			}
			if result {
				fmt.Println("Login successful")
				client.Login(name)
				client.Write(messagingProtocol.NewSuccess(""))
			} else {
				fmt.Println("Login failed")
				client.Write(messagingProtocol.NewError("Login failed"))
			}
		case protocol.CREATE_EVENT:
			fmt.Println("user wants to create an event")

			if len(data.Data) < 3 {
				fmt.Println("Invalid number of arguments")
				client.Write(messagingProtocol.NewError("Invalid number of arguments"))
				continue
			}
			if client.state != connected {
				fmt.Println("User is not logged in")
				client.Write(messagingProtocol.NewError("You must be logged in to create an event"))
				continue
			}

			eventName := data.Data[0]

			_, err := db.CreateEvent(eventName, client.GetConnected())
			if err != nil {
				fmt.Println("Error creating event: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			}
			fmt.Println("Event created")
			event, err := db.GetEventByName(eventName)
			if err != nil {
				fmt.Println("Error creating event: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			}
			for i := 1; i < len(data.Data)-1; i += 2 {
				nbVolunteers, err := strconv.ParseUint(data.Data[i+1], 10, 32)
				if err != nil {
					fmt.Println("Error parsing number of volunteers: ", err.Error())
					client.Write(messagingProtocol.NewError(err.Error()))
					continue
				}
				event.CreateJob(data.Data[i], uint(nbVolunteers))
			}
			fmt.Println("Jobs created")
			client.Write(messagingProtocol.NewSuccess(""))
			client.Logout()

		case protocol.GET_EVENTS:
			fmt.Println("user wants to get events")

			err := client.Write(protocol.DataPacket{
				Type: protocol.OK,
				Data: db.GetEventsAsStringArray(),
			})

			if err != nil {
				fmt.Println("Error sending events: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			} else {
				fmt.Println("Events sent")
				for _, event := range db.GetEventsAsStringArray() {
					fmt.Println("Event: ", event)
				}
			}
		case protocol.GET_JOBS:
			fmt.Println("user wants to get jobs")

			if len(data.Data) != 1 {
				fmt.Println("Invalid number of arguments")
				client.Write(messagingProtocol.NewError("Invalid number of arguments"))
				continue
			}

			eventName := data.Data[0]

			event, err := db.GetEventByName(eventName)
			if err != nil {
				fmt.Println("Error getting event: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			}

			err = client.Write(protocol.DataPacket{
				Type: protocol.OK,
				Data: event.GetJobsAsStringArray(),
			})

			if err != nil {
				fmt.Println("Error sending jobs: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			} else {
				fmt.Println("Jobs sent")
				for _, job := range event.GetJobsAsStringArray() {
					fmt.Println("Job: ", job)
				}
			}

		case protocol.JOIN_EVENT:
			fmt.Println("user wants to join an event")

			if len(data.Data) != 2 {
				fmt.Println("Invalid number of arguments")
				client.Write(messagingProtocol.NewError("Invalid number of arguments"))
				continue
			}

			if client.state != connected {
				fmt.Println("User is not logged in")
				client.Write(messagingProtocol.NewError("You must be logged in to join an event"))
				continue
			}

			eventName := data.Data[0]
			jobName := data.Data[1]

			event, err := db.GetEventByName(eventName)
			if err != nil {
				fmt.Println("Error getting event: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			}

			job, err := event.GetJob(jobName)
			if err != nil {
				fmt.Println("Error getting job: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			}

			_, err = job.AddVolunteer(client.GetConnected())
			if err != nil {
				fmt.Println("Error adding volunteer: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			}

			fmt.Println("Volunteer added")
			client.Write(messagingProtocol.NewSuccess(""))
			client.Logout()

		case protocol.STOP:
			fmt.Println("user wants to stop the server")
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}
