package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Database is a in memory structure holding event and user data
type Database struct {
	Events map[uint]*Event
	Users  map[string]*User
}

// jsonDatabase is an helper structure used to serialize/deserialize the database to/from json.
type jsonDatabase struct {
	Events jsonEvents `json:"events"`
	Users  Users      `json:"users"`
}

// LoadDatabaseFromJson returns the state of the database stored in the json file at path
func LoadDatabaseFromJson(jsonPath string) Database {
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var conf jsonDatabase
	json.Unmarshal(byteValue, &conf)

	return Database{conf.Events.ToMap(), conf.Users.ToMap()}
}

// GetEventByName returns the event with the name `name`
//
// Returns an error if no such event exists
// Complexity O(n)
func (db *Database) GetEventByName(name string) (*Event, error) {
	for _, event := range db.Events {
		if event.Name == name {
			return event, nil
		}
	}
	return nil, ErrorEventNotFound
}

// GetEvent returns the event with the id `id`
//
// Returns an error if no such event exists
// Complexity O(1)
func (db *Database) GetEvent(id uint) (*Event, error) {
	if event, ok := db.Events[id]; ok {
		return event, nil
	}
	return nil, ErrorEventNotFound
}

// GetEventAsStringArray returns the string representation of all the events in the database
func (db *Database) GetEventsAsStringArray() []string {
	var events []string
	for _, event := range db.Events {
		events = append(events, event.ToString())
	}
	return events
}

// CreateEvent creates an event in the database
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
	id := uint(len(db.Events))
	db.Events[id] = &Event{id, name, organizer, make(map[uint]*Job), true}
	return db, nil
}

// CreateUser creates a new user in the database
func (db *Database) CreateUser(name string, password string, function string) (*Database, error) {
	if name == "" {
		return nil, ErrorUserNameEmpty
	}
	if password == "" {
		return nil, ErrorPasswordEmpty
	}
	if function == "" {
		return nil, ErrorFunctionEmpty
	}
	if _, err := db.GetUser(name); err == nil {
		return nil, ErrorUserExists
	}
	usr := User{name, password, function}
	db.Users[name] = &usr
	return db, nil
}

// GetUser returns the user with the given name
//
// Complexity: O(1)
func (db *Database) GetUser(name string) (*User, error) {
	if name == "" {
		return &User{}, ErrorUserNameEmpty
	}
	if user, ok := db.Users[name]; ok {
		return user, nil
	}
	return &User{}, ErrorUserNotFound
}

// Login checks if the user exists and if the password is correct
func (db Database) Login(name string, password string) (bool, error) {
	user, err := db.GetUser(name)
	if err != nil {
		return false, err
	}
	return user.Password == password, nil
}
