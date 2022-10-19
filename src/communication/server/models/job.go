package models

import "errors"

type Job struct {
	Name       string   `json:"name"`
	Required   int      `json:"required"`
	Volunteers []string `json:"volunteers"`
}

// Creates a new job and adds it to the database
// Returns an error if a job with the same id already exists
// Otherwise returns the new state of the database
func (event Event) CreateJob(name string, required int) (Event, error) {
	if _, err := event.GetJob(name); err == nil {
		return event, errors.New("Job with same id already exists")
	}
	event.Jobs = append(event.Jobs, Job{name, required, []string{}})
	return event, nil
}

// Get a job from the database
// Returns an error if the job does not exist
// Otherwise returns the job
func (event Event) GetJob(name string) (Job, error) {
	for _, job := range event.Jobs {
		if job.Name == name {
			return job, nil
		}
	}
	return Job{}, errors.New("Job not found")
}

// Adds a volunteer to a job
// Returns an error if the job does not exist
// Otherwise returns the new state of the database
func (job Job) AddVolunteer(userId string) (Job, error) {
	job.Volunteers = append(job.Volunteers, userId)
	return job, nil
}
