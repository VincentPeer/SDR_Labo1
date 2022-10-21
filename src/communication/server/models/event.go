package models

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrorEventNotFound  = errors.New("event not found")
	ErrorEventExists    = errors.New("event with same id already exists")
	ErrorEventNameEmpty = errors.New("event name cannot be empty")
	ErrorOrganizerEmpty = errors.New("organizer name cannot be empty")
	ErrorNotOrganizer   = errors.New("user is not an organizer")
	ErrorEventIsClosed  = errors.New("event is closed")
)

// event holds the events' data
type event struct {
	ID        uint
	Name      string
	Organizer string
	Jobs      map[uint]*job
	isOpen    bool
}

// jsonEvent is an helper structure used to serialize/deserialize the event to/from json.
type jsonEvent struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Organizer string `json:"organizer"`
	Jobs      jobs   `json:"jobs"`
	IsOpen    bool   `json:"isOpen"`
}

type jsonEvents []jsonEvent

// toMap converts the json structure to a map of events
func (e *jsonEvents) toMap() map[uint]*event {
	events := make(map[uint]*event)
	for i := 0; i < len(*e); i++ {
		events[(*e)[i].ID] = &event{(*e)[i].ID, (*e)[i].Name, (*e)[i].Organizer, (*e)[i].Jobs.ToMap(), (*e)[i].IsOpen}
	}
	return events
}

// CreateJob creates a new job in the database
func (e *event) CreateJob(name string, required uint) (*event, error) {
	if name == "" {
		return nil, ErrorJobNameEmpty
	}
	if _, err := e.GetJobByName(name); err == nil {
		return e, ErrorJobExists
	}
	id := uint(len(e.Jobs))
	e.Jobs[id] = &job{ID: id, Name: name, Required: required, Volunteers: []string{}}
	return e, nil
}

// GetJob returns the job with the given id
//
// Complexity: O(1)
func (e *event) GetJob(id uint) (*job, error) {
	job, found := e.Jobs[id]
	if !found {
		return nil, ErrorJobNotFound
	}
	return job, nil
}

// GetJobByName returns the job with the given name
//
// Complexity: O(n)
func (e *event) GetJobByName(name string) (*job, error) {
	if name == "" {
		return nil, ErrorJobNameEmpty
	}
	for _, job := range e.Jobs {
		if job.Name == name {
			return job, nil
		}
	}
	return &job{}, ErrorJobNotFound
}

// ToString returns a string representation of the event
func (e *event) ToString() string {
	openStatus := "open"
	if !e.isOpen {
		openStatus = "closed"
	}
	return fmt.Sprintf("%d | %s | %s | %s", e.ID, e.Name, e.Organizer, openStatus)
}

// GetJobAsStringArray returns the jobs as an array of strings
func (e *event) GetJobsAsStringArray() []string {
	var jobs []string
	for _, job := range e.Jobs {
		jobs = append(jobs, job.ToString())
	}
	return jobs
}

// GEtJobsRepartitionTable returns a table with the jobs and which volunteers are assigned to them
func (e *event) GetJobsRepartitionTable() []string {
	var table []string
	for _, job := range e.Jobs {
		line := fmt.Sprintf(" %d : ", job.ID)
		for _, volunteer := range job.Volunteers {
			line += fmt.Sprintf("%s - ", volunteer)
		}
		table = append(table, line)
	}
	return table
}

func (e *event) GetJobsRepartitionTable2() []string {
	head := "| Volunteers | "
	for _, job := range e.Jobs {
		s := fmt.Sprintf("%-15s", job.Name+" "+strconv.Itoa((int)(job.Required))+" |")
		head += s
	}
	var tab []string
	tab = append(tab, head)
	for _, job := range e.Jobs {
		for _, volunteer := range job.Volunteers {
			s := fmt.Sprintf("%-10s", volunteer)
			tab = append(tab, "| "+s+" | ")
		}
	}
	return tab
}

// AddVolunteerToJob adds a volunteer to a job
//
// If the volunteer is already assigned to a job in the event, it is removed from that job
func (e *event) AddVolunteer(jobId uint, name string) (*job, error) {
	if name == "" {
		return nil, ErrorVolunteerEmpty
	}
	job, err := e.GetJob(jobId)
	if err != nil {
		return nil, err
	}
	if !e.isOpen {
		return nil, ErrorEventIsClosed
	}
	if job.Required == uint(len(job.Volunteers)) {
		return nil, ErrorVolunteerAboveMaximum
	}
	if _, err := job.GetVolunteer(name); err != ErrorVolunteerNotFound {
		return job, ErrorVolunteerExists
	}
	err = e.RemoveVolunteer(name)
	if err != nil {
		return nil, err
	}
	job.Volunteers = append(job.Volunteers, name)
	return job, nil
}

// RemoveVolunteer removes a volunteer from the jobs in the event
func (e *event) RemoveVolunteer(name string) error {
	if name == "" {
		return ErrorVolunteerEmpty
	}
	for _, job := range e.Jobs {
		if _, err := job.GetVolunteer(name); err == nil {
			job.RemoveVolunteer(name)
		}
	}
	return nil
}

// Close closes the event
//
// This means that no more volunteers can be added to the event
func (e *event) Close() {
	e.isOpen = false
}
