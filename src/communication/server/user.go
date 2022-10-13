package main

import "errors"

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Function string `json:"function"`
}

// Creates a new user and adds it to the database
// Returns an error if a user with the same id already exists
// Otherwise returns the new state of the database
func createUser(db []User, id string, name string, password string, function string) ([]User, error) {
	if _, err := getUser(db, id); err == nil {
		return db, errors.New("User with same id already exists")
	}
	return append(db, User{id, name, password, function}), nil
}

// Get a user from the database
// Returns an error if the user does not exist
// Otherwise returns the user
func getUser(db []User, id string) (User, error) {
	for _, user := range db {
		if user.Id == id {
			return user, nil
		}
	}
	return User{}, errors.New("User not found")
}
