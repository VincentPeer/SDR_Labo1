package models

import (
	"errors"
	"fmt"
)

var (
	ErrorJobNotFound           = errors.New("j not found")
	ErrorJobExists             = errors.New("j with same id already exists")
	ErrorJobNameEmpty          = errors.New("j name cannot be empty")
	ErrorVolunteerNotFound     = errors.New("volunteer not found")
	ErrorVolunteerExists       = errors.New("volunteer already signed up for this j")
	ErrorVolunteerEmpty        = errors.New("volunteer name cannot be empty")
	ErrorVolunteerAboveMaximum = errors.New("volunteer count is above maximum")
)

type jobs []job

// job holds the jobs' data
type job struct {
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	Required   uint     `json:"required"`
	Volunteers []string `json:"volunteers"`
}

// ToMap converts the json structure to a map of jobs
func (js *jobs) ToMap() map[uint]*job {
	jobsMap := make(map[uint]*job)
	for i := 0; i < len(*js); i++ {
		jobsMap[(*js)[i].ID] = &(*js)[i]
	}
	return jobsMap
}

// GetVolunteer returns the volunteer with the given name
//
// Complexity: O(n)
func (j *job) GetVolunteer(name string) (string, error) {
	if name == "" {
		return "", ErrorVolunteerEmpty
	}
	for _, volunteer := range j.Volunteers {
		if volunteer == name {
			return volunteer, nil
		}
	}
	return "", ErrorVolunteerNotFound
}

// RemoveVolunteer removes the volunteer with the given name
func (j *job) RemoveVolunteer(name string) (*job, error) {
	if name == "" {
		return nil, ErrorVolunteerEmpty
	}
	for i, volunteer := range j.Volunteers {
		if volunteer == name {
			j.Volunteers = append(j.Volunteers[:i], j.Volunteers[i+1:]...)
			return j, nil
		}
	}
	return nil, ErrorVolunteerNotFound
}

// ToString converts the j to a string
func (j *job) ToString() string {
	return fmt.Sprintf("%d | %-10s | %d |", j.ID, j.Name, j.Required)
}
