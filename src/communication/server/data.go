package main

type Config struct {
	Users  []User  `json:"users"`
	Events []Event `json:"events"`
}

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Function string `json:"function"`
}

type Event struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Organizer User   `json:"organizer"`
	Jobs      []Job  `json:"jobs"`
}

type Job struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Required   int    `json:"required"`
	Volunteers []User `json:"volunteers"`
}
