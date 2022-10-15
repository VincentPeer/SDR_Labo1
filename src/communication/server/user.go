package main

import "errors"

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Function string `json:"function"`
}

// Creates a new user and adds it to the database
// Returns an error if a user with the same id already exists
// Otherwise returns the new state of the database
func createUser(db []User, name string, password string, function string) ([]User, error) {
	if _, err := getUser(db, name); err == nil {
		return nil, errors.New("User with same id already exists")
	}
	return append(db, User{name, password, function}), nil
}

// Get a user from the database
// Returns an error if the user does not exist
// Otherwise returns the user
func getUser(db []User, name string) (User, error) {
	for _, user := range db {
		if user.Name == name {
			return user, nil
		}
	}
	return User{}, errors.New("User not found")
}

// Confirm the password of a user
// Returns an error if the user does not exist
func login(name string, password string) (bool, error) {
	user, err := getUser(users, name)
	if err != nil {
		return false, err
	}
	return user.Password == password, nil
}
