package server

import (
	"SDR_Labo1/src/communication/protocol"
	"SDR_Labo1/src/communication/server/models"
	"strconv"
	"time"
)

// DatabaseManager handles all database requests and updates.
//
// It listens to a channel for DatabaseRequest and handles them. This way, only one thread can access them
type DatabaseManager struct {
	db             models.Database
	RequestChannel chan DatabaseRequest
	isDebug        bool
}

// DatabaseRequest represents an incoming database request
//
// Client (goroutine) who originated the message, to use in responses
type DatabaseRequest struct {
	sender  *clientConnection
	payload protocol.DataPacket
}

// NewDatabaseRequest returns a database request initialised with payload as input and the goroutine origin
func NewDatabaseRequest(client *clientConnection, data protocol.DataPacket) *DatabaseRequest {
	return &DatabaseRequest{
		sender:  client,
		payload: data,
	}
}

// NewDatabaseManager returns a database manager and initiate its channel
func NewDatabaseManager(db models.Database, isDebug bool) *DatabaseManager {
	return &DatabaseManager{
		db:             db,
		RequestChannel: make(chan DatabaseRequest),
		isDebug:        isDebug,
	}
}

// IsDebug return true if debugging prints should be performed
func (dbm *DatabaseManager) IsDebug() bool {
	return dbm.isDebug
}

// Start starts the thread, listening to the user request channel. This is a blocking function.
func (dbm *DatabaseManager) Start() {
	for {
		request := <-dbm.RequestChannel
		dbm.handleRequest(request)
	}
}

// handleRequest handles a request by parsing it and performing the requested action
func (dbm *DatabaseManager) handleRequest(request DatabaseRequest) {
	if request.sender.isDebug && dbm.isDebug {
		time.Sleep(5 * time.Second)
	}

	switch request.payload.Type {
	case protocol.LOGIN:
		Debug(dbm, "user wants to login")
		loginHandler(dbm, request)
	case protocol.CREATE_EVENT:
		Debug(dbm, "user wants to create an event")
		createEventHandler(dbm, request)
	case protocol.GET_EVENTS:
		Debug(dbm, "user wants to get events")
		getEventsHandler(dbm, request)
	case protocol.GET_JOBS:
		Debug(dbm, "user wants to get jobs")
		getJobsHandler(dbm, request)
	case protocol.EVENT_REG:
		Debug(dbm, "user wants to join an event")
		eventRegHandler(dbm, request)
	case protocol.CLOSE_EVENT:
		Debug(dbm, "user wants to close an event")
		closeEventHandler(dbm, request)
	case protocol.STOP:
		Debug(dbm, "user wants to stop the a")
		stopHandler(dbm, request)
	default:
		Debug(dbm, "Unknown command")
	}
	if request.payload.Type != protocol.LOGIN && request.payload.Type != protocol.STOP {
		request.sender.connectedUser = ""
	}
	Debug(dbm, "Request handled")
}

// checkDatapacket checks the number of parameters of a request
//
// returns true if everything is ok, false otherwise
func checkDatapacket(data protocol.DataPacket, minNbParams int, maxNbParams int, client *clientConnection) bool {
	if len(data.Data) < minNbParams || len(data.Data) > maxNbParams {
		client.sendError("Invalid number of arguments")
		return false
	}
	return true
}

// logInUser checks if a user can login with the given credentials
//
// if login is successfull the client connected user is updated. An update tcp message is also sent to the client
func (dbm *DatabaseManager) logInUser(client *clientConnection, username string, password string) (bool, error) {
	Debug(dbm, "user wants to login")
	Debug(dbm, "name: "+username)
	Debug(dbm, " password: "+password)

	user, err := dbm.db.GetUser(username)
	errMsg := "Login failed"
	if err != nil {
		errMsg = err.Error()
	} else if user.Password == password {
		client.connectedUser = username
		client.sendSuccess("Login successful")
		return true, nil
	}
	Debug(dbm, errMsg)
	client.sendError(errMsg)
	return false, err
}

// checkIfConnected checks that the client is connected
func checkIfConnected(client *clientConnection) bool {
	if client.connectedUser == "" {
		client.sendError("You must be logged in to do this")
		return false
	}
	return true
}

// checkIfOrganizer checks if the connected user to a given client is the organizer of a given event
func checkIfOrganizer(client *clientConnection, event *models.Event) bool {
	if event.Organizer != client.connectedUser {
		client.sendError("You are not the organizer of this event")
		return false
	}
	return true
}

// parseInt parses a string that should represents an int
func parseInt(client *clientConnection, data string) (int, error) {
	integer, err := strconv.ParseInt(data, 10, 32)
	if err != nil {
		client.sendError("Invalid integer")
		return 0, err
	}
	return int(integer), nil
}
