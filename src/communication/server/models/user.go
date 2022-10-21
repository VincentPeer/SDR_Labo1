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

// user holds the users' data
type user struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Function string `json:"function"`
}

type users []user

// ToMap converts the json structure to a map of users
func (us *users) ToMap() map[string]*user {
	usersMap := make(map[string]*user)
	for i := range *us {
		usersMap[(*us)[i].Name] = &(*us)[i]
	}
	return usersMap
}
