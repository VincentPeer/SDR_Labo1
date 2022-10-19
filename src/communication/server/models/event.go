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
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Organizer string `json:"organizer"`
	Jobs      []Job  `json:"jobs"`
}

type Events []Event

// Creates a new job and adds it to the database
// Returns an error if a job with the same id already exists
// Otherwise returns the new state of the database
func (event *Event) CreateJob(name string, required uint) (*Event, error) {
	if name == "" {
		return nil, ErrorJobNameEmpty
	}
	if _, err := event.GetJob(name); err == nil {
		return event, ErrorJobExists
	}

	event.Jobs = append(event.Jobs, Job{uint(len(event.Jobs)), name, required, []string{}})
	return event, nil
}

// Get a job from the database
// Returns an error if the job does not exist
// Otherwise returns the job
func (event *Event) GetJob(name string) (*Job, error) {
	if name == "" {
		return nil, ErrorJobNameEmpty
	}
	for _, job := range event.Jobs {
		if job.Name == name {
			return &job, nil
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
