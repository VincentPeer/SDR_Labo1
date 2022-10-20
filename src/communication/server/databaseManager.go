package server

import (
	"SDR_Labo1/src/communication/protocol"
	"SDR_Labo1/src/communication/server/models"
	"fmt"
	"math"
	"strconv"
)

type DatabaseManager struct {
	db             models.Database
	RequestChannel chan DatabaseRequest
}

type DatabaseRequest struct {
	sender  *Client
	payload protocol.DataPacket
}

func NewDatabaseRequest(client *Client, data protocol.DataPacket) *DatabaseRequest {
	return &DatabaseRequest{
		sender:  client,
		payload: data,
	}
}

func NewDatabaseManager(db models.Database) *DatabaseManager {
	return &DatabaseManager{
		db:             db,
		RequestChannel: make(chan DatabaseRequest),
	}
}

func (dbm *DatabaseManager) Start() {
	for {
		request := <-dbm.RequestChannel
		dbm.handleRequest(request)
	}
}

func (dbm *DatabaseManager) handleRequest(request DatabaseRequest) {
	switch request.payload.Type {
	case protocol.LOGIN:
		if checkDatapacket(request.payload, 2, 2, request.sender) {
			dbm.logInUser(request.sender, request.payload.Data[0], request.payload.Data[1])
		}
	case protocol.CREATE_EVENT:
		fmt.Println("user wants to create an event")
		if checkDatapacket(request.payload, 1, math.MaxInt32, request.sender) && checkIfConnected(request.sender) {
			eventName := request.payload.Data[0]
			_, err := dbm.db.CreateEvent(eventName, request.sender.GetConnected())
			if err != nil {
				fmt.Println(err.Error())
				request.sender.SendError(err.Error())
				break
			}
			event, err := dbm.db.GetEventByName(eventName)
			if err != nil {
				fmt.Println(err.Error())
				request.sender.SendError(err.Error())
				break
			}
			// Populating the event with jobs
			for i := 1; i < len(request.payload.Data)-1; i += 2 {
				nbVolunteers, err := strconv.ParseUint(request.payload.Data[i+1], 10, 32)
				if err != nil {
					fmt.Println("Error parsing number of volunteers: ", err.Error())
					request.sender.SendError(err.Error())
					break
				}
				event.CreateJob(request.payload.Data[i], uint(nbVolunteers))
			}
			request.sender.SendSuccess("Event created")
			fmt.Println("Jobs created")
		}
	case protocol.GET_EVENTS:
		fmt.Println("user wants to get events")

		if len(request.payload.Data) == 0 { // GET all events
			err := request.sender.Write(protocol.DataPacket{
				Type: protocol.OK,
				Data: dbm.db.GetEventsAsStringArray(),
			})
			if err != nil {
				fmt.Println("Error sending events: ", err.Error())
				request.sender.SendError(err.Error())
				break
			}
			fmt.Println("Events sent")

		} else if len(request.payload.Data) == 1 { // GET all jobs for an event
			eventId, err := strconv.ParseUint(request.payload.Data[0], 10, 32)
			if err != nil {
				fmt.Println("Invalid eventId: ", request.payload.Data[0])
				request.sender.SendError("Invalid eventId: is not a uint64")
				break
			}
			event, err := dbm.db.GetEvent(uint(eventId))
			if err != nil {
				request.sender.SendError(err.Error())
				fmt.Println("Error getting event: ", err.Error())
				break
			}
			err = request.sender.Write(protocol.DataPacket{
				Type: protocol.OK,
				Data: event.GetJobsAsStringArray(),
			})
			if err != nil {
				fmt.Println("Error getting events: ", err.Error())
				request.sender.SendError(err.Error())
				break
			}
			fmt.Println("events sent")
		} else {
			fmt.Println("ERROR: wrong number of arguments")
			request.sender.SendError("Incorrect number of arguments.\nNeed 0 or 1 (eventID)")
		}
	case protocol.GET_JOBS:
		fmt.Println("user wants to get jobs")

		if checkDatapacket(request.payload, 1, 1, request.sender) {
			eventId, err := strconv.ParseUint(request.payload.Data[0], 10, 32)
			if err != nil {
				fmt.Println("Invalid eventId: ", request.payload.Data[0])
				request.sender.SendError("Invalid eventId: is not a uint64")
				break
			}
			event, err := dbm.db.GetEvent(uint(eventId))
			if err != nil {
				request.sender.SendError(err.Error())
				fmt.Println("Error getting event: ", err.Error())
				break
			}
			err = request.sender.Write(protocol.DataPacket{
				Type: protocol.OK,
				Data: event.GetJobsAsStringArray(),
			})
			if err != nil {
				fmt.Println("Error sending jobs: ", err.Error())
				request.sender.SendError(err.Error())
				break
			}
			fmt.Println("events sent")
		}

	case protocol.EVENT_REG:
		fmt.Println("user wants to join an event")

		if checkDatapacket(request.payload, 2, 2, request.sender) && checkIfConnected(request.sender) {
			eventId, err := parseInt(request.sender, request.payload.Data[0])
			if err != nil {
				break
			}
			jobId, err := parseInt(request.sender, request.payload.Data[1])
			if err != nil {
				break
			}
			event, err := dbm.db.GetEvent(uint(eventId))
			if err != nil {
				fmt.Println("Error getting event: ", err.Error())
				request.sender.SendError(err.Error())
				break
			}

			job, err := event.GetJob(uint(jobId))
			if err != nil {
				fmt.Println("Error getting job: ", err.Error())
				request.sender.SendError(err.Error())
				break
			}
			fmt.Println(event.GetJobsRepartitionTable())
			fmt.Println(job)

			_, err = event.AddVolunteer(job.ID, request.sender.GetConnected())
			if err != nil {
				fmt.Println("Error adding volunteer: ", err.Error())
				request.sender.SendError(err.Error())
				break
			}

			fmt.Println("Volunteer added")
			fmt.Println(event.GetJobsRepartitionTable())
			fmt.Println(job)
			request.sender.SendSuccess("Volunteer added")
		}
	case protocol.CLOSE_EVENT:
		fmt.Println("user wants to close an event")

		if checkDatapacket(request.payload, 1, 1, request.sender) && checkIfConnected(request.sender) {

			eventId, err := parseInt(request.sender, request.payload.Data[0])
			if err != nil {
				break
			}
			event, err := dbm.db.GetEvent(uint(eventId))
			if err != nil {
				fmt.Println("Error getting event: ", err.Error())
				request.sender.SendError(err.Error())
				break
			}

			if checkIfOrganizer(request.sender, event) {
				event.Close()
				request.sender.SendSuccess("Event closed")
			}
		}

	case protocol.STOP:
		fmt.Println("user wants to stop the a")
		request.sender.Close()
		return

	default:
		fmt.Println("Unknown command")
	}
	if request.payload.Type != protocol.LOGIN {
		request.sender.Logout()
	}
	fmt.Println("Request handled")
}

func checkDatapacket(data protocol.DataPacket, minNbParams int, maxNbParams int, client *Client) bool {
	if len(data.Data) < minNbParams || len(data.Data) > maxNbParams {
		fmt.Println("Invalid number of arguments")
		client.SendError("Invalid number of arguments")
		return false
	}
	return true
}

func (dbm *DatabaseManager) logInUser(client *Client, username string, password string) (bool, error) {
	fmt.Println("user wants to login")
	fmt.Print("name: ", username)
	fmt.Println(" password: ", password)

	user, err := dbm.db.GetUser(username)
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

func checkIfConnected(client *Client) bool {
	if client.state != connected {
		fmt.Println("User is not logged in")
		client.SendError("You must be logged in to do this")
		return false
	}
	return true
}

func checkIfOrganizer(client *Client, event *models.Event) bool {
	if event.Organizer != client.GetConnected() {
		fmt.Println("User is not the organizer")
		client.SendError("You are not the organizer of this event")
		return false
	}
	return true
}

func parseInt(client *Client, data string) (int, error) {
	integer, err := strconv.ParseInt(data, 10, 32)
	if err != nil {
		fmt.Println("Invalid integer: ", data)
		client.SendError("Invalid integer")
		return 0, err
	}
	return int(integer), nil
}
