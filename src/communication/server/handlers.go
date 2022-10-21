package server

import (
	"SDR_Labo1/src/communication/protocol"
	"math"
	"strconv"
	"strings"
)

func loginHandler(dbm *DatabaseManager, request DatabaseRequest) {
	if checkDatapacket(request.payload, 2, 2, request.sender) {
		dbm.logInUser(request.sender, request.payload.Data[0], request.payload.Data[1])
	}
}

func createEventHandler(dbm *DatabaseManager, request DatabaseRequest) {
	if checkDatapacket(request.payload, 1, math.MaxInt32, request.sender) && checkIfConnected(request.sender) {
		eventName := request.payload.Data[0]
		_, err := dbm.db.CreateEvent(eventName, request.sender.connectedUser)
		if err != nil {
			Debug(dbm, err.Error())
			request.sender.sendError(err.Error())
			return
		}
		event, err := dbm.db.GetEventByName(eventName)
		if err != nil {
			Debug(dbm, err.Error())
			request.sender.sendError(err.Error())
			return
		}
		// Populating the event with jobs
		for i := 1; i < len(request.payload.Data)-1; i += 2 {
			nbVolunteers, err := strconv.ParseUint(request.payload.Data[i+1], 10, 32)
			if err != nil {
				Debug(dbm, "Error parsing number of volunteers: "+err.Error())
				request.sender.sendError(err.Error())
				return
			}
			event.CreateJob(request.payload.Data[i], uint(nbVolunteers))
		}
		request.sender.sendSuccess("Event created")
		Debug(dbm, "Event created")
	}
}

func getEventsHandler(dbm *DatabaseManager, request DatabaseRequest) {

	if len(request.payload.Data) == 0 { // GET all events
		err := request.sender.write(protocol.DataPacket{
			Type: protocol.OK,
			Data: dbm.db.GetEventsAsStringArray(),
		})
		if err != nil {
			Debug(dbm, "Error sending events: "+err.Error())
			request.sender.sendError(err.Error())
			return
		}
		Debug(dbm, "Events sent")

	} else if len(request.payload.Data) == 1 { // GET all jobs for an event
		eventId, err := strconv.ParseUint(request.payload.Data[0], 10, 32)
		if err != nil {
			Debug(dbm, "Invalid eventId: "+request.payload.Data[0])
			request.sender.sendError("Invalid eventId: is not a uint64")
			return
		}
		event, err := dbm.db.GetEvent(uint(eventId))
		if err != nil {
			request.sender.sendError(err.Error())
			Debug(dbm, "Error getting event: "+err.Error())
			return
		}
		err = request.sender.write(protocol.DataPacket{
			Type: protocol.OK,
			Data: event.GetJobsAsStringArray(),
		})
		if err != nil {
			Debug(dbm, "Error getting events: "+err.Error())
			request.sender.sendError(err.Error())
			return
		}
		Debug(dbm, "events sent")
	} else {
		Debug(dbm, "ERROR: wrong number of arguments")
		request.sender.sendError("Incorrect number of arguments.\nNeed 0 or 1 (eventID)")
	}
}

func getJobsHandler(dbm *DatabaseManager, request DatabaseRequest) {
	if checkDatapacket(request.payload, 1, 1, request.sender) {
		eventId, err := strconv.ParseUint(request.payload.Data[0], 10, 32)
		if err != nil {
			Debug(dbm, "Invalid eventId: "+request.payload.Data[0])
			request.sender.sendError("Invalid eventId: is not a uint64")
			return
		}
		event, err := dbm.db.GetEvent(uint(eventId))
		if err != nil {
			request.sender.sendError(err.Error())
			Debug(dbm, "Error getting event: "+err.Error())
			return
		}
		err = request.sender.write(protocol.DataPacket{
			Type: protocol.OK,
			Data: event.GetJobsRepartitionTable2(),
		})
		if err != nil {
			Debug(dbm, "Error sending jobs: "+err.Error())
			request.sender.sendError(err.Error())
			return
		}
		Debug(dbm, "events sent")
	}
}

func eventRegHandler(dbm *DatabaseManager, request DatabaseRequest) {
	if checkDatapacket(request.payload, 2, 2, request.sender) && checkIfConnected(request.sender) {
		eventId, err := parseInt(request.sender, request.payload.Data[0])
		if err != nil {
			return
		}
		jobId, err := parseInt(request.sender, request.payload.Data[1])
		if err != nil {
			return
		}
		event, err := dbm.db.GetEvent(uint(eventId))
		if err != nil {
			Debug(dbm, "Error getting event: "+err.Error())
			request.sender.sendError(err.Error())
			return
		}

		job, err := event.GetJob(uint(jobId))
		if err != nil {
			Debug(dbm, "Error getting job: "+err.Error())
			request.sender.sendError(err.Error())
			return
		}
		Debug(dbm, strings.Join(event.GetJobsRepartitionTable(), "\n"))
		Debug(dbm, job.ToString())

		_, err = event.AddVolunteer(job.ID, request.sender.connectedUser)
		if err != nil {
			Debug(dbm, "Error adding volunteer: "+err.Error())
			request.sender.sendError(err.Error())
			return
		}

		Debug(dbm, "Volunteer added")
		Debug(dbm, strings.Join(event.GetJobsRepartitionTable(), "\n"))
		Debug(dbm, job.ToString())
		request.sender.sendSuccess("Volunteer added")
	}
}

func closeEventHandler(dbm *DatabaseManager, request DatabaseRequest) {
	if checkDatapacket(request.payload, 1, 1, request.sender) && checkIfConnected(request.sender) {

		eventId, err := parseInt(request.sender, request.payload.Data[0])
		if err != nil {
			return
		}
		event, err := dbm.db.GetEvent(uint(eventId))
		if err != nil {
			Debug(dbm, "Error getting event: "+err.Error())
			request.sender.sendError(err.Error())
			return
		}

		if checkIfOrganizer(request.sender, event) {
			event.Close()
			request.sender.sendSuccess("Event closed")
		}
	}
}

func stopHandler(dbm *DatabaseManager, request DatabaseRequest) {
	request.sender.close()
}
