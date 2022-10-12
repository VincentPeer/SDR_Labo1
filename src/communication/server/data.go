package main

type Config struct {
	Users  []User  `json:"users"`
	Events []Event `json:"events"`
}

type Event struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Organizer string `json:"organizer"`
	Jobs      []Job  `json:"jobs"`
}

type Job struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Required   int      `json:"required"`
	Volunteers []string `json:"volunteers"`
}
