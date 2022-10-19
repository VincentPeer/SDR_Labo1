package models

import "errors"

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Function string `json:"function"`
}

type Users []User

// Creates a new user and adds it to the database
// Returns an error if a user with the same id already exists
// Otherwise returns the new state of the database
func (db Users) CreateUser(name string, password string, function string) (Users, error) {
	if _, err := db.GetUser(name); err == nil {
		return nil, errors.New("User with same id already exists")
	}
	return append(db, User{name, password, function}), nil
}

// Get a user from the database
// Returns an error if the user does not exist
// Otherwise returns the user
func (db Users) GetUser(name string) (User, error) {
	for _, user := range db {
		if user.Name == name {
			return user, nil
		}
	}
	return User{}, errors.New("User not found")
}

// Confirm the password of a user
// Returns an error if the user does not exist
func (db Users) Login(name string, password string) (bool, error) {
	user, err := db.GetUser(name)
	if err != nil {
		return false, err
	}
	return user.Password == password, nil
}
