package models

import "errors"

type Event struct {
	Name      string `json:"name"`
	Organizer string `json:"organizer"`
	Jobs      []Job  `json:"jobs"`
}

type Events []Event

// Creates a new event and adds it to the database
// Returns an error if an event with the same id already exists
// Otherwise returns the new state of the database
func (db Events) CreateEvent(name string, organizer string) (Events, error) {
	if _, err := db.GetEvent(name); err == nil {
		return nil, errors.New("Event with same id already exists")
	}
	return append(db, Event{name, organizer, []Job{}}), nil
}

// Get an event from the database
// Returns an error if the event does not exist
// Otherwise returns the event
func (db Events) GetEvent(name string) (Event, error) {
	for _, event := range db {
		if event.Name == name {
			return event, nil
		}
	}
	return Event{}, errors.New("Event not found")
}
