package main

import "errors"

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Function string `json:"function"`
}

// Create a new user
func createUser(id string, name string, password string, function string) {
	users = append(users, User{id, name, password, function})
}

func getUser(id string) (User, error) {
	for _, user := range users {
		if user.Id == id {
			return user, nil
		}
	}
	return User{}, errors.New("User not found")
}
