package models

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrorJobNotFound           = errors.New("job not found")
	ErrorJobExists             = errors.New("job with same id already exists")
	ErrorJobNameEmpty          = errors.New("job name cannot be empty")
	ErrorVolunteerNotFound     = errors.New("volunteer not found")
	ErrorVolunteerExists       = errors.New("volunteer already signed up for this job")
	ErrorVolunteerEmpty        = errors.New("volunteer name cannot be empty")
	ErrorVolunteerAboveMaximum = errors.New("volunteer count is above maximum")
)

type Job struct {
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	Required   uint     `json:"required"`
	Volunteers []string `json:"volunteers"`
	EventId    uint     `json:"event_id"`
}

// Get a volunteer from the database
// Returns an error if the volunteer does not exist
func (job *Job) GetVolunteer(name string) (string, error) {
	if name == "" {
		return "", ErrorVolunteerEmpty
	}
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
	if name == "" {
		return nil, ErrorVolunteerEmpty
	}
	if job.Required == uint(len(job.Volunteers)) {
		return nil, ErrorVolunteerAboveMaximum
	}
	if _, err := job.GetVolunteer(name); err != ErrorVolunteerNotFound {
		return job, ErrorVolunteerExists
	}
	job.Volunteers = append(job.Volunteers, name)
	return job, nil
}

func (job *Job) ToString() string {
	return fmt.Sprintf("%d | %s | %d | %s", job.ID, job.Name, job.Required, strings.Join(job.Volunteers, " - "))
}
