package models

import (
	"errors"
	"fmt"
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

type Jobs []Job

type Job struct {
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	Required   uint     `json:"required"`
	Volunteers []string `json:"volunteers"`
	EventId    uint     `json:"event_id"`
}

func (jobs *Jobs) ToMap() map[uint]*Job {
	jobsMap := make(map[uint]*Job)
	for i := 0; i < len(*jobs); i++ {
		jobsMap[(*jobs)[i].ID] = &(*jobs)[i]
	}
	return jobsMap
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

// Remove a volunteer from the database
// Returns an error if the volunteer does not exist
// Otherwise returns the new state of the database
func (job *Job) RemoveVolunteer(name string) (*Job, error) {
	if name == "" {
		return nil, ErrorVolunteerEmpty
	}
	for i, volunteer := range job.Volunteers {
		if volunteer == name {
			job.Volunteers = append(job.Volunteers[:i], job.Volunteers[i+1:]...)
			return job, nil
		}
	}
	return nil, ErrorVolunteerNotFound
}

func (job *Job) ToString() string {
	return fmt.Sprintf("%d | %s | %d", job.ID, job.Name, job.Required)
}
