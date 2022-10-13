package main

type Config struct {
	Users  []User  `json:"users"`
	Events []Event `json:"events"`
}
