package main

var Event struct {
	eventName  string
	organizer  Person
	posts      []Post // pointeur?
	volunteers []Person
}
