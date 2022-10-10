package event

import (
	"bufio"
	"fmt"
)

type Event struct { // todo est-ce que exporter est bien? (pas poo)
	eventName    string
	organizer    User
	posts        []Post // pointeur?
	volunteers   []User
	nbVolunteers uint
}

func CreateEvent(readWriter *bufio.ReadWriter) {
	var event Event

	fmt.Println("Enter the event's name : ")
	event.eventName, _ = readWriter.ReadString('\n')

	fmt.Println("Enter your username : ")
	event.organizer.name, _ = readWriter.ReadString('\n')

	fmt.Println("Enter your password : ")
	event.organizer.password, _ = readWriter.ReadString('\n')

	fmt.Println("Enter the each post's name, termin with double enter : ")
	i := 0
	for { // todo voir comment se rempli le tableau non alloué
		name, _ := readWriter.ReadString('\n')
		if name == "\n" {
			break
		}
		event.posts[i].name = name
		i++
	}

	fmt.Println("Enter the number of volunteer for this event : ")
	//event.nbVolunteers, _ = readWriter.ReadString('\n')
}
