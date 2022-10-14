package main

import "errors"

type Event struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Organizer string `json:"organizer"`
	Jobs      []Job  `json:"jobs"`
}

type Job struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Required   int      `json:"required"`
	Volunteers []string `json:"volunteers"`
}

// Creates a new event and adds it to the database
// Returns an error if an event with the same id already exists
// Otherwise returns the new state of the database
func createEvent(db []Event, id string, name string, organizer string) ([]Event, error) {
	if _, err := getEvent(db, id); err == nil {
		return nil, errors.New("Event with same id already exists")
	}
	return append(db, Event{id, name, organizer, []Job{}}), nil
}

// Get an event from the database
// Returns an error if the event does not exist
// Otherwise returns the event
func getEvent(db []Event, id string) (Event, error) {
	for _, event := range db {
		if event.Id == id {
			return event, nil
		}
	}
	return Event{}, errors.New("Event not found")
}

// Creates a new job and adds it to the database
// Returns an error if a job with the same id already exists
// Otherwise returns the new state of the database
func createJob(event Event, id string, name string, required int) (Event, error) {
	if _, err := getJob(event, id); err == nil {
		return event, errors.New("Job with same id already exists")
	}
	event.Jobs = append(event.Jobs, Job{id, name, required, []string{}})
	return event, nil
}

// Get a job from the database
// Returns an error if the job does not exist
// Otherwise returns the job
func getJob(event Event, id string) (Job, error) {
	for _, job := range event.Jobs {
		if job.Id == id {
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
