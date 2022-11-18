package client

// Database is a in memory structure holding event and user data
type Database struct {
	Events map[uint]*event
}

// jsonDatabase is an helper structure used to serialize/deserialize the database to/from json.
type jsonDatabase struct {
	Events jsonEvents
}

// event holds the events' data
type event struct {
	ID        uint
	Name      string
	Organizer string
	Jobs      map[uint]*job
	isOpen    bool
}

// jsonEvent is an helper structure used to serialize/deserialize the event to/from json.
type jsonEvent struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Organizer string `json:"organizer"`
	Jobs      jobs   `json:"jobs"`
	IsOpen    bool   `json:"isOpen"`
}

type jsonEvents []jsonEvent

// toMap converts the json structure to a map of events
func (e *jsonEvents) toMap() map[uint]*event {
	events := make(map[uint]*event)
	for i := 0; i < len(*e); i++ {
		events[(*e)[i].ID] = &event{(*e)[i].ID, (*e)[i].Name, (*e)[i].Organizer, (*e)[i].Jobs.ToMap(), (*e)[i].IsOpen}
	}
	return events
}

type jobs []job

// job holds the jobs' data
type job struct {
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	Required   uint     `json:"required"`
	Volunteers []string `json:"volunteers"`
}

// ToMap converts the json structure to a map of jobs
func (js *jobs) ToMap() map[uint]*job {
	jobsMap := make(map[uint]*job)
	for i := 0; i < len(*js); i++ {
		jobsMap[(*js)[i].ID] = &(*js)[i]
	}
	return jobsMap
}

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
