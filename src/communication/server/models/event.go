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
)

type Event struct {
	ID        uint
	Name      string
	Organizer string
	Jobs      map[uint]*Job
	isOpen    bool
}

type jsonEvent struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Organizer string `json:"organizer"`
	Jobs      Jobs   `json:"jobs"`
	isOpen    bool   `json:"isOpen"`
}

type jsonEvents []jsonEvent
type Events []Event

func (event *jsonEvents) ToMap() map[uint]*Event {
	events := make(map[uint]*Event)
	for i := 0; i < len(*event); i++ {
		events[(*event)[i].ID] = &Event{(*event)[i].ID, (*event)[i].Name, (*event)[i].Organizer, (*event)[i].Jobs.ToMap(), (*event)[i].isOpen}
	}
	return events
}

// Creates a new job and adds it to the database
// Returns an error if a job with the same id already exists
// Otherwise returns the new state of the database
func (event *Event) CreateJob(name string, required uint) (*Event, error) {
	if name == "" {
		return nil, ErrorJobNameEmpty
	}
	if _, err := event.GetJobByName(name); err == nil {
		return event, ErrorJobExists
	}
	id := uint(len(event.Jobs))
	event.Jobs[id] = &Job{ID: id, Name: name, Required: required, Volunteers: []string{}, EventId: event.ID}
	return event, nil
}

// Get a job from the database
// Returns an error if the job does not exist
// Otherwise returns the job
func (event *Event) GetJob(id uint) (*Job, error) {
	//	if id == "" {
	//		return nil, ErrorJobNameEmpty
	//	}
	job, found := event.Jobs[id]
	if !found {
		return nil, ErrorJobNotFound
	}
	return job, nil
}

// Get a job from the database
// Returns an error if the job does not exist
// Otherwise returns the job
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

func (event *Event) ToString() string {
	return fmt.Sprintf("%d | %s | %s", event.ID, event.Name, event.Organizer)
}

func (event *Event) GetJobsAsStringArray() []string {
	var jobs []string
	for _, job := range event.Jobs {
		jobs = append(jobs, job.ToString())
	}
	return jobs
}

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

// Adds a volunteer to a job
// Returns an error if the job does not exist
// Otherwise returns the new state of the database
func (event *Event) AddVolunteer(jobId uint, name string) (*Job, error) {
	job, err := event.GetJob(jobId)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, ErrorVolunteerEmpty
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

// Removes a volunteer from all jobs
// Returns an error if the volunteer does not exist
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
