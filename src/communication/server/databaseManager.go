package server

import (
	"SDR_Labo1/src/communication/protocol"
	"SDR_Labo1/src/communication/server/models"
	"math"
	"strconv"
	"strings"
	"time"
)

type DatabaseManager struct {
	db             models.Database
	RequestChannel chan DatabaseRequest
	isDebug        bool
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

func NewDatabaseManager(db models.Database, isDebug bool) *DatabaseManager {
	return &DatabaseManager{
		db:             db,
		RequestChannel: make(chan DatabaseRequest),
		isDebug:        isDebug,
	}
}

func (dbm *DatabaseManager) IsDebug() bool {
	return dbm.isDebug
}

func (dbm *DatabaseManager) Start() {
	for {
		request := <-dbm.RequestChannel
		dbm.handleRequest(request)
	}
}

func (dbm *DatabaseManager) handleRequest(request DatabaseRequest) {
	if request.sender.isDebug && dbm.isDebug {
		time.Sleep(5 * time.Second)
	}

	switch request.payload.Type {
	case protocol.LOGIN:
		if checkDatapacket(request.payload, 2, 2, request.sender) {
			dbm.logInUser(request.sender, request.payload.Data[0], request.payload.Data[1])
		}
	case protocol.CREATE_EVENT:
		Debug(dbm, "user wants to create an event")
		if checkDatapacket(request.payload, 1, math.MaxInt32, request.sender) && checkIfConnected(request.sender) {
			eventName := request.payload.Data[0]
			_, err := dbm.db.CreateEvent(eventName, request.sender.GetConnected())
			if err != nil {
				Debug(dbm, err.Error())
				request.sender.SendError(err.Error())
				break
			}
			event, err := dbm.db.GetEventByName(eventName)
			if err != nil {
				Debug(dbm, err.Error())
				request.sender.SendError(err.Error())
				break
			}
			// Populating the event with jobs
			for i := 1; i < len(request.payload.Data)-1; i += 2 {
				nbVolunteers, err := strconv.ParseUint(request.payload.Data[i+1], 10, 32)
				if err != nil {
					Debug(dbm, "Error parsing number of volunteers: "+err.Error())
					request.sender.SendError(err.Error())
					break
				}
				event.CreateJob(request.payload.Data[i], uint(nbVolunteers))
			}
			request.sender.SendSuccess("Event created")
			Debug(dbm, "Event created")
		}
	case protocol.GET_EVENTS:
		Debug(dbm, "user wants to get events")

		if len(request.payload.Data) == 0 { // GET all events
			err := request.sender.Write(protocol.DataPacket{
				Type: protocol.OK,
				Data: dbm.db.GetEventsAsStringArray(),
			})
			if err != nil {
				Debug(dbm, "Error sending events: "+err.Error())
				request.sender.SendError(err.Error())
				break
			}
			Debug(dbm, "Events sent")

		} else if len(request.payload.Data) == 1 { // GET all jobs for an event
			eventId, err := strconv.ParseUint(request.payload.Data[0], 10, 32)
			if err != nil {
				Debug(dbm, "Invalid eventId: "+request.payload.Data[0])
				request.sender.SendError("Invalid eventId: is not a uint64")
				break
			}
			event, err := dbm.db.GetEvent(uint(eventId))
			if err != nil {
				request.sender.SendError(err.Error())
				Debug(dbm, "Error getting event: "+err.Error())
				break
			}
			err = request.sender.Write(protocol.DataPacket{
				Type: protocol.OK,
				Data: event.GetJobsAsStringArray(),
			})
			if err != nil {
				Debug(dbm, "Error getting events: "+err.Error())
				request.sender.SendError(err.Error())
				break
			}
			Debug(dbm, "events sent")
		} else {
			Debug(dbm, "ERROR: wrong number of arguments")
			request.sender.SendError("Incorrect number of arguments.\nNeed 0 or 1 (eventID)")
		}
	case protocol.GET_JOBS:
		Debug(dbm, "user wants to get jobs")

		if checkDatapacket(request.payload, 1, 1, request.sender) {
			eventId, err := strconv.ParseUint(request.payload.Data[0], 10, 32)
			if err != nil {
				Debug(dbm, "Invalid eventId: "+request.payload.Data[0])
				request.sender.SendError("Invalid eventId: is not a uint64")
				break
			}
			event, err := dbm.db.GetEvent(uint(eventId))
			if err != nil {
				request.sender.SendError(err.Error())
				Debug(dbm, "Error getting event: "+err.Error())
				break
			}
			err = request.sender.Write(protocol.DataPacket{
				Type: protocol.OK,
				Data: event.GetJobsAsStringArray(),
			})
			if err != nil {
				Debug(dbm, "Error sending jobs: "+err.Error())
				request.sender.SendError(err.Error())
				break
			}
			Debug(dbm, "events sent")
		}

	case protocol.EVENT_REG:
		Debug(dbm, "user wants to join an event")

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
				Debug(dbm, "Error getting event: "+err.Error())
				request.sender.SendError(err.Error())
				break
			}

			job, err := event.GetJob(uint(jobId))
			if err != nil {
				Debug(dbm, "Error getting job: "+err.Error())
				request.sender.SendError(err.Error())
				break
			}
			Debug(dbm, strings.Join(event.GetJobsRepartitionTable(), "\n"))
			Debug(dbm, job.ToString())

			_, err = event.AddVolunteer(job.ID, request.sender.GetConnected())
			if err != nil {
				Debug(dbm, "Error adding volunteer: "+err.Error())
				request.sender.SendError(err.Error())
				break
			}

			Debug(dbm, "Volunteer added")
			Debug(dbm, strings.Join(event.GetJobsRepartitionTable(), "\n"))
			Debug(dbm, job.ToString())
			request.sender.SendSuccess("Volunteer added")
		}
	case protocol.CLOSE_EVENT:
		Debug(dbm, "user wants to close an event")

		if checkDatapacket(request.payload, 1, 1, request.sender) && checkIfConnected(request.sender) {

			eventId, err := parseInt(request.sender, request.payload.Data[0])
			if err != nil {
				break
			}
			event, err := dbm.db.GetEvent(uint(eventId))
			if err != nil {
				Debug(dbm, "Error getting event: "+err.Error())
				request.sender.SendError(err.Error())
				break
			}

			if checkIfOrganizer(request.sender, event) {
				event.Close()
				request.sender.SendSuccess("Event closed")
			}
		}

	case protocol.STOP:
		Debug(dbm, "user wants to stop the a")
		request.sender.Close()
		return

	default:
		Debug(dbm, "Unknown command")
	}
	if request.payload.Type != protocol.LOGIN {
		request.sender.Logout()
	}
	Debug(dbm, "Request handled")
}

func checkDatapacket(data protocol.DataPacket, minNbParams int, maxNbParams int, client *Client) bool {
	if len(data.Data) < minNbParams || len(data.Data) > maxNbParams {
		client.SendError("Invalid number of arguments")
		return false
	}
	return true
}

func (dbm *DatabaseManager) logInUser(client *Client, username string, password string) (bool, error) {
	Debug(dbm, "user wants to login")
	Debug(dbm, "name: "+username)
	Debug(dbm, " password: "+password)

	user, err := dbm.db.GetUser(username)
	errMsg := "Login failed"
	if err != nil {
		errMsg = err.Error()
	} else if user.Password == password {
		client.Login(username)
		client.SendSuccess("Login successful")
		return true, nil
	}
	Debug(dbm, errMsg)
	client.SendError(errMsg)
	return false, err
}

func checkIfConnected(client *Client) bool {
	if client.state != connected {
		client.SendError("You must be logged in to do this")
		return false
	}
	return true
}

func checkIfOrganizer(client *Client, event *models.Event) bool {
	if event.Organizer != client.GetConnected() {
		client.SendError("You are not the organizer of this event")
		return false
	}
	return true
}

func parseInt(client *Client, data string) (int, error) {
	integer, err := strconv.ParseInt(data, 10, 32)
	if err != nil {
		client.SendError("Invalid integer")
		return 0, err
	}
	return int(integer), nil
}
