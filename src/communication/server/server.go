package server

import (
	"SDR_Labo1/src/communication/protocol"
	"SDR_Labo1/src/communication/server/models"
	"fmt"
	"io"
	"math"
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

type Server struct {
	clients           map[int]*client
	db                models.Database
	messagingProtocol protocol.TcpProtocol
}

func NewServer() *Server {
	path, err := filepath.Abs(CONFIG_FILE_PATH)

	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}

	return &Server{
		clients:           make(map[int]*client),
		db:                models.LoadDatabaseFromJson(path),
		messagingProtocol: protocol.TcpProtocol{},
	}
}

func (server *Server) Start() {

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

func closeRequest(client *client) {
	fmt.Println("Closing client")
	client.Close()
}

// Handles incoming requests.
func (server *Server) handleRequest(client *client) {
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
			if checkDatapacket(data, 2, 2, client) {
				logInUser(client, data.Data[0], data.Data[1])
			}
		case protocol.CREATE_EVENT:
			fmt.Println("user wants to create an event")
			if checkDatapacket(data, 1, math.MaxInt32, client) && checkIfConnected(client) {
				eventName := data.Data[0]
				_, err := client.server.db.CreateEvent(eventName, client.GetConnected())
				if err != nil {
					fmt.Println(err.Error())
					client.SendError(err.Error())
					continue
				}
				event, err := server.db.GetEventByName(eventName)
				if err != nil {
					fmt.Println(err.Error())
					client.SendError(err.Error())
					continue
				}
				// Populating the event with jobs
				for i := 1; i < len(data.Data)-1; i += 2 {
					nbVolunteers, err := strconv.ParseUint(data.Data[i+1], 10, 32)
					if err != nil {
						fmt.Println("Error parsing number of volunteers: ", err.Error())
						client.SendError(err.Error())
						continue
					}
					event.CreateJob(data.Data[i], uint(nbVolunteers))
				}
				client.SendSuccess("Event created")
				fmt.Println("Jobs created")
			}
		case protocol.GET_EVENTS:
			fmt.Println("user wants to get events")

			if len(data.Data) == 0 { // GET all events
				err := client.Write(protocol.DataPacket{
					Type: protocol.OK,
					Data: server.db.GetEventsAsStringArray(),
				})
				if err != nil {
					fmt.Println("Error sending events: ", err.Error())
					client.SendError(err.Error())
					continue
				}
				fmt.Println("Events sent")

			} else if len(data.Data) == 1 { // GET all jobs for an event
				eventId, err := strconv.ParseUint(data.Data[0], 10, 32)
				if err != nil {
					fmt.Println("Invalid eventId: ", data.Data[0])
					client.SendError("Invalid eventId: is not a uint64")
					continue
				}
				event, err := server.db.GetEvent(uint(eventId))
				if err != nil {
					client.SendError(err.Error())
					fmt.Println("Error getting event: ", err.Error())
					continue
				}
				err = client.Write(protocol.DataPacket{
					Type: protocol.OK,
					Data: event.GetJobsAsStringArray(),
				})
				if err != nil {
					fmt.Println("Error getting events: ", err.Error())
					client.SendError(err.Error())
					continue
				}
				fmt.Println("events sent")
			} else {
				fmt.Println("ERROR: wrong number of arguments")
				client.SendError("Incorrect number of arguments.\nNeed 0 or 1 (eventID)")
			}
		case protocol.GET_JOBS:
			fmt.Println("user wants to get jobs")

			if checkDatapacket(data, 1, 1, client) {
				eventId, err := strconv.ParseUint(data.Data[0], 10, 32)
				if err != nil {
					fmt.Println("Invalid eventId: ", data.Data[0])
					client.SendError("Invalid eventId: is not a uint64")
					continue
				}
				event, err := server.db.GetEvent(uint(eventId))
				if err != nil {
					client.SendError(err.Error())
					fmt.Println("Error getting event: ", err.Error())
					continue
				}
				err = client.Write(protocol.DataPacket{
					Type: protocol.OK,
					Data: event.GetJobsAsStringArray(),
				})
				if err != nil {
					fmt.Println("Error sending jobs: ", err.Error())
					client.SendError(err.Error())
					continue
				}
				fmt.Println("events sent")
			}

		case protocol.EVENT_REG:
			fmt.Println("user wants to join an event")

			if checkDatapacket(data, 2, 2, client) && checkIfConnected(client) {
				eventId, err := parseInt(client, data.Data[0])
				if err != nil {
					continue
				}
				jobId, err := parseInt(client, data.Data[1])
				if err != nil {
					continue
				}
				event, err := server.db.GetEvent(uint(eventId))
				if err != nil {
					fmt.Println("Error getting event: ", err.Error())
					client.SendError(err.Error())
					continue
				}

				job, err := event.GetJob(uint(jobId))
				if err != nil {
					fmt.Println("Error getting job: ", err.Error())
					client.SendError(err.Error())
					continue
				}
				fmt.Println(event.GetJobsRepartitionTable())
				fmt.Println(job)

				_, err = event.AddVolunteer(job.ID, client.GetConnected())
				if err != nil {
					fmt.Println("Error adding volunteer: ", err.Error())
					client.SendError(err.Error())
					continue
				}

				fmt.Println("Volunteer added")
				fmt.Println(event.GetJobsRepartitionTable())
				fmt.Println(job)
				client.SendSuccess("Volunteer added")
			}
		case protocol.CLOSE_EVENT:
			fmt.Println("user wants to close an event")

			if checkDatapacket(data, 1, 1, client) && checkIfConnected(client) {

				eventId, err := parseInt(client, data.Data[0])
				if err != nil {
					continue
				}
				event, err := server.db.GetEvent(uint(eventId))
				if err != nil {
					fmt.Println("Error getting event: ", err.Error())
					client.SendError(err.Error())
					continue
				}

				if !checkIfOrganizer(client, event) {
					event.Close()
					client.SendSuccess("Event closed")
				}
			}

		case protocol.STOP:
			fmt.Println("user wants to stop the a")
			return

		default:
			fmt.Println("Unknown command")
		}
		client.Logout()
	}
}

func checkDatapacket(data protocol.DataPacket, minNbParams int, maxNbParams int, client *client) bool {
	if len(data.Data) < minNbParams || len(data.Data) > maxNbParams {
		fmt.Println("Invalid number of arguments")
		client.SendError("Invalid number of arguments")
		return false
	}
	return true
}

func logInUser(client *client, username string, password string) (bool, error) {
	fmt.Println("user wants to login")
	fmt.Print("name: ", username)
	fmt.Println(" password: ", password)

	user, err := client.server.db.GetUser(username)
	errMsg := "Login failed"
	if err != nil {
		errMsg = err.Error()
	} else if user.Password == password {
		client.Login(username)
		client.SendSuccess("Login successful")
		return true, nil
	}
	fmt.Println(errMsg)
	client.SendError(errMsg)
	return false, err
}

func checkIfConnected(client *client) bool {
	if client.state != connected {
		fmt.Println("User is not logged in")
		client.SendError("You must be logged in to do this")
		return false
	}
	return true
}

func checkIfOrganizer(client *client, event *models.Event) bool {
	if event.Organizer != client.GetConnected() {
		fmt.Println("User is not the organizer")
		client.SendError("You are not the organizer of this event")
		return false
	}
	return true
}

func parseInt(client *client, data string) (int, error) {
	integer, err := strconv.ParseInt(data, 10, 32)
	if err != nil {
		fmt.Println("Invalid integer: ", data)
		client.SendError("Invalid integer")
		return 0, err
	}
	return int(integer), nil
}
