package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Database struct {
	Events Events `json:"events"`
	Users  Users  `json:"users"`
}

func LoadDatabaseFromJson(jsonPath string) Database {
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}
	//fmt.Println("Successfully opened " + jsonFile.Name())
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var conf Database
	json.Unmarshal(byteValue, &conf)

	return conf
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

func (db *Database) GetEventsAsStringArray() []string {
	var events []string
	for _, event := range db.Events {
		events = append(events, event.ToString())
	}
	return events
}

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
	db.Events = append(db.Events, Event{uint(len(db.Events)), name, organizer, []Job{}})
	return db, nil
}
