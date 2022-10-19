package models

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrorJobNotFound       = errors.New("job not found")
	ErrorJobExists         = errors.New("job with same id already exists")
	ErrorJobNameEmpty      = errors.New("job name cannot be empty")
	ErrorVolunteerNotFound = errors.New("volunteer not found")
	ErrorVolunteerExists   = errors.New("volunteer already signed up for this job")
)

type Job struct {
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	Required   uint     `json:"required"`
	Volunteers []string `json:"volunteers"`
}

// Creates a new job and adds it to the database
// Returns an error if a job with the same id already exists
// Otherwise returns the new state of the database
func (event *Event) CreateJob(name string, required uint) (*Event, error) {
	if name == "" {
		return nil, ErrorJobNameEmpty
	}
	if _, err := event.GetJob(name); err == nil {
		return event, errors.New("Job with same id already exists")
	}

	event.Jobs = append(event.Jobs, Job{uint(len(event.Jobs)), name, required, []string{}})
	return event, nil
}

// Get a job from the database
// Returns an error if the job does not exist
// Otherwise returns the job
func (event *Event) GetJob(name string) (*Job, error) {
	for _, job := range event.Jobs {
		if job.Name == name {
			return &job, nil
		}
	}
	return &Job{}, errors.New("Job not found")
}

// Get a volunteer from the database
// Returns an error if the volunteer does not exist
func (job *Job) GetVolunteer(name string) (string, error) {
	for _, volunteer := range job.Volunteers {
		if volunteer == name {
			return volunteer, nil
		}
	}
	return "", ErrorVolunteerNotFound
}

// Adds a volunteer to a job
// Returns an error if the job does not exist
// Otherwise returns the new state of the database
func (job *Job) AddVolunteer(name string) (*Job, error) {
	if _, err := job.GetVolunteer(name); err != ErrorVolunteerNotFound {
		return job, ErrorVolunteerExists
	}
	job.Volunteers = append(job.Volunteers, name)
	return job, nil
}

func (job *Job) ToString() string {
	return fmt.Sprintf("%d | %s | %d | %s", job.ID, job.Name, job.Required, strings.Join(job.Volunteers, " - "))
}
