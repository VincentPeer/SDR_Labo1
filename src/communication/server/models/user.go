package models

import (
	"errors"
)

var (
	ErrorUserNotFound  = errors.New("user not found")
	ErrorUserExists    = errors.New("user with same id already exists")
	ErrorUserNameEmpty = errors.New("user name cannot be empty")
	ErrorPasswordEmpty = errors.New("password cannot be empty")
	ErrorFunctionEmpty = errors.New("function cannot be empty")
)

// User holds the users' data
type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Function string `json:"function"`
}

type Users []User

// ToMap converts the json structure to a map of users
func (users *Users) ToMap() map[string]*User {
	usersMap := make(map[string]*User)
	for i, _ := range *users {
		usersMap[(*users)[i].Name] = &(*users)[i]
	}
	return usersMap
}
