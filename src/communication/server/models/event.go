package models

import (
	"errors"
	"fmt"
)

var (
	ErrorEventNotFound  = errors.New("event not found")
	ErrorEventExists    = errors.New("event with same id already exists")
	ErrorEventNameEmpty = errors.New("event name cannot be empty")
	ErrorOrganizerEmpty = errors.New("organizer name cannot be empty")
	ErrorNotOrganizer   = errors.New("user is not an organizer")
	ErrorEventIsClosed  = errors.New("event is closed")
)

// Event holds the events' data
type Event struct {
	ID        uint
	Name      string
	Organizer string
	Jobs      map[uint]*Job
	isOpen    bool
}

// jsonEvent is an helper structure used to serialize/deserialize the event to/from json.
type jsonEvent struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Organizer string `json:"organizer"`
	Jobs      Jobs   `json:"jobs"`
	IsOpen    bool   `json:"isOpen"`
}

type jsonEvents []jsonEvent
type Events []Event

// ToMap converts the json structure to a map of events
func (event *jsonEvents) ToMap() map[uint]*Event {
	events := make(map[uint]*Event)
	for i := 0; i < len(*event); i++ {
		events[(*event)[i].ID] = &Event{(*event)[i].ID, (*event)[i].Name, (*event)[i].Organizer, (*event)[i].Jobs.ToMap(), (*event)[i].IsOpen}
	}
	return events
}

// CreateJob creates a new job in the database
func (event *Event) CreateJob(name string, required uint) (*Event, error) {
	if name == "" {
		return nil, ErrorJobNameEmpty
	}
	if _, err := event.GetJobByName(name); err == nil {
		return event, ErrorJobExists
	}
	id := uint(len(event.Jobs))
	event.Jobs[id] = &Job{ID: id, Name: name, Required: required, Volunteers: []string{}}
	return event, nil
}

// GetJob returns the job with the given id
//
// Complexity: O(1)
func (event *Event) GetJob(id uint) (*Job, error) {
	job, found := event.Jobs[id]
	if !found {
		return nil, ErrorJobNotFound
	}
	return job, nil
}

// GetJobByName returns the job with the given name
//
// Complexity: O(n)
func (event *Event) GetJobByName(name string) (*Job, error) {
	if name == "" {
		return nil, ErrorJobNameEmpty
	}
	for _, job := range event.Jobs {
		if job.Name == name {
			return job, nil
		}
	}
	return &Job{}, ErrorJobNotFound
}

// ToString returns a string representation of the event
func (event *Event) ToString() string {
	openStatus := "open"
	if !event.isOpen {
		openStatus = "closed"
	}
	return fmt.Sprintf("%d | %s | %s | %s", event.ID, event.Name, event.Organizer, openStatus)
}

// GetJobAsStringArray returns the jobs as an array of strings
func (event *Event) GetJobsAsStringArray() []string {
	var jobs []string
	for _, job := range event.Jobs {
		jobs = append(jobs, job.ToString())
	}
	return jobs
}

// GEtJobsRepartitionTable returns a table with the jobs and which volunteers are assigned to them
func (event *Event) GetJobsRepartitionTable() []string {
	var table []string
	for _, job := range event.Jobs {
		line := fmt.Sprintf(" %d : ", job.ID)
		for _, volunteer := range job.Volunteers {
			line += fmt.Sprintf("%s - ", volunteer)
		}
		table = append(table, line)
	}
	return table
}

// AddVolunteerToJob adds a volunteer to a job
//
// If the volunteer is already assigned to a job in the event, it is removed from that job
func (event *Event) AddVolunteer(jobId uint, name string) (*Job, error) {
	if name == "" {
		return nil, ErrorVolunteerEmpty
	}
	job, err := event.GetJob(jobId)
	if err != nil {
		return nil, err
	}
	if !event.isOpen {
		return nil, ErrorEventIsClosed
	}
	if job.Required == uint(len(job.Volunteers)) {
		return nil, ErrorVolunteerAboveMaximum
	}
	if _, err := job.GetVolunteer(name); err != ErrorVolunteerNotFound {
		return job, ErrorVolunteerExists
	}
	err = event.RemoveVolunteer(name)
	if err != nil {
		return nil, err
	}
	job.Volunteers = append(job.Volunteers, name)
	return job, nil
}

// RemoveVolunteer removes a volunteer from the jobs in the event
func (event *Event) RemoveVolunteer(name string) error {
	if name == "" {
		return ErrorVolunteerEmpty
	}
	for _, job := range event.Jobs {
		if _, err := job.GetVolunteer(name); err == nil {
			job.RemoveVolunteer(name)
		}
	}
	return nil
}

// Close closes the event
//
// This means that no more volunteers can be added to the event
func (event *Event) Close() {
	event.isOpen = false
}
