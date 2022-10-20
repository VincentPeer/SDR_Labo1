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
			for _, user := range db.Users {
				fmt.Println("username : ", user.Name)
				fmt.Println("userpass : ", user.Password)
			}
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

			if len(data.Data) == 0 { // GET all events

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
			} else if len(data.Data) == 1 { // GET all jobs for an event
				eventId, err := strconv.ParseUint(data.Data[0], 10, 32)
				if err != nil {
					fmt.Println("Invalid eventId: ", data.Data[0])
					client.Write(messagingProtocol.NewError("Invalid eventId: is not a uint64"))
					continue
				}

				event, err := db.GetEvent(uint(eventId))
				if err != nil {
					client.Write(messagingProtocol.NewError(err.Error()))
					fmt.Println("Error getting event: ", err.Error())
					continue
				}
				err = client.Write(protocol.DataPacket{
					Type: protocol.OK,
					Data: event.GetJobsAsStringArray(),
				})

				if err != nil {
					fmt.Println("Error getting events: ", err.Error())
					client.Write(messagingProtocol.NewError(err.Error()))
					continue
				} else {
					fmt.Println("events sent")
					for _, jobID := range event.GetJobsAsStringArray() {
						fmt.Println("eventID: ", jobID)
					}
				}
			} else {
				fmt.Println("ERROR: wrong number of arguments")
				client.Write(messagingProtocol.NewError("Incorrect number of arguments.\nNeed 0 or 1 (eventID)"))
			}

		case protocol.GET_JOBS:
			fmt.Println("user wants to get jobs")

			if len(data.Data) != 1 {
				fmt.Println("Invalid number of arguments")
				client.Write(messagingProtocol.NewError("Invalid number of arguments"))
				continue
			}

			eventId, err := strconv.ParseUint(data.Data[0], 10, 32)
			if err != nil {
				fmt.Println("Invalid eventId: ", data.Data[0])
				client.Write(messagingProtocol.NewError("Invalid eventId: is not a uint64"))
				continue
			}

			event, err := db.GetEvent(uint(eventId))
			if err != nil {
				fmt.Println("Error getting event: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			}

			err = client.Write(protocol.DataPacket{
				Type: protocol.OK,
				Data: event.GetJobsRepartitionTable(),
			})

			if err != nil {
				fmt.Println("Error sending jobs: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			} else {
				fmt.Println("Jobs sent")
				for _, job := range event.GetJobsRepartitionTable() {
					fmt.Println("Job: ", job)
				}
			}

		case protocol.EVENT_REG:
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

			eventId, err := strconv.ParseUint(data.Data[0], 10, 32)
			if err != nil {
				client.Write(messagingProtocol.NewError("Argument 1 (eventId) is not a uint64"))
				fmt.Println(err.Error())
				continue
			}
			jobId, err := strconv.ParseUint(data.Data[1], 10, 32)
			if err != nil {
				fmt.Println("Invalid eventId: ", data.Data[0])
				client.Write(messagingProtocol.NewError("Invalid eventId: is not a uint64"))
				continue
			}
			event, err := db.GetEvent(uint(eventId))
			if err != nil {
				fmt.Println("Error getting event: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			}

			job, err := event.GetJob(uint(jobId))
			if err != nil {
				fmt.Println("Error getting job: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			}
			fmt.Println(event.GetJobsRepartitionTable())
			fmt.Println(job)
			_, err = event.AddVolunteer(job.ID, client.GetConnected())
			if err != nil {
				fmt.Println("Error adding volunteer: ", err.Error())
				client.Write(messagingProtocol.NewError(err.Error()))
				continue
			}

			fmt.Println("Volunteer added")
			fmt.Println(event.GetJobsRepartitionTable())
			fmt.Println(job)
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
