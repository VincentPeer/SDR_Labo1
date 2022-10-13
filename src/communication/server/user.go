package main

import "errors"

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Function string `json:"function"`
}

// Create a new user
func createUser(db []User, id string, name string, password string, function string) []User {
	return append(db, User{id, name, password, function})
}

func getUser(db []User, id string) (User, error) {
	for _, user := range db {
		if user.Id == id {
			return user, nil
		}
	}
	return User{}, errors.New("User not found")
}
