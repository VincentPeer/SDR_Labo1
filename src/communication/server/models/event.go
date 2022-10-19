package models

import "errors"

var (
	ErrorEventNotFound  = errors.New("event not found")
	ErrorEventExists    = errors.New("event with same id already exists")
	ErrorEventNameEmpty = errors.New("event name cannot be empty")
	ErrorOrganizerEmpty = errors.New("organizer name cannot be empty")
	ErrorNotOrganizer   = errors.New("user is not an organizer")
)

type Event struct {
	Name      string `json:"name"`
	Organizer string `json:"organizer"`
	Jobs      []Job  `json:"jobs"`
}

type Events []Event

// Creates a new event and adds it to the database
// Returns an error if an event with the same id already exists
// Otherwise returns the new state of the database
func (db *Database) CreateEvent(name string, organizer string) (*Database, error) {
	if name == "" {
		return nil, ErrorEventNameEmpty
	}
	if organizer == "" {
		return nil, ErrorOrganizerEmpty
	}
	user, err := db.GetUser(organizer)
	if err != nil {
		return nil, err
	}
	if user.Function != "organizer" {
		return nil, ErrorNotOrganizer
	}

	if _, err := db.GetEvent(name); err == nil {
		return db, ErrorEventExists
	}
	db.Events = append(db.Events, Event{name, organizer, []Job{}})
	return db, nil
}

// Get an event from the database
// Returns an error if the event does not exist
// Otherwise returns the event
func (db *Database) GetEvent(name string) (Event, error) {
	for _, event := range db.Events {
		if event.Name == name {
			return event, nil
		}
	}
	return Event{}, ErrorEventNotFound
}
