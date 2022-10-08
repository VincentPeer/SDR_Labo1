package event

var Event struct {
	eventName  string
	organizer  User
	posts      []Post // pointeur?
	volunteers []User
}
