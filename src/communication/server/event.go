package main

import "errors"

type Event struct {
	Name      string `json:"name"`
	Organizer string `json:"organizer"`
	Jobs      []Job  `json:"jobs"`
}

type Job struct {
	Name       string   `json:"name"`
	Required   int      `json:"required"`
	Volunteers []string `json:"volunteers"`
}

// Creates a new event and adds it to the database
// Returns an error if an event with the same id already exists
// Otherwise returns the new state of the database
func createEvent(db []Event, name string, organizer string) ([]Event, error) {
	if _, err := getEvent(db, name); err == nil {
		return nil, errors.New("Event with same id already exists")
	}
	return append(db, Event{name, organizer, []Job{}}), nil
}

// Get an event from the database
// Returns an error if the event does not exist
// Otherwise returns the event
func getEvent(db []Event, name string) (Event, error) {
	for _, event := range db {
		if event.Name == name {
			return event, nil
		}
	}
	return Event{}, errors.New("Event not found")
}

// Creates a new job and adds it to the database
// Returns an error if a job with the same id already exists
// Otherwise returns the new state of the database
func createJob(event Event, name string, required int) (Event, error) {
	if _, err := getJob(event, name); err == nil {
		return event, errors.New("Job with same id already exists")
	}
	event.Jobs = append(event.Jobs, Job{name, required, []string{}})
	return event, nil
}

// Get a job from the database
// Returns an error if the job does not exist
// Otherwise returns the job
func getJob(event Event, name string) (Job, error) {
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
func addVolunteer(job Job, userId string) (Job, error) {
	job.Volunteers = append(job.Volunteers, userId)
	return job, nil
}
