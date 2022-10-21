package client

import (
	"SDR_Labo1/src/communication/protocol"
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var messagingProtocol = &protocol.SDRProtocol{}
var consoleIn = bufio.NewReader(os.Stdin)

// userInterface is the main function for the user interface,
// the client can go through each different functionality
func userInterface(c *Connection) {
	fmt.Println("Welcome!")

	var choice int
	for {
		fmt.Println("Choose one of the following functionality")
		fmt.Println("[1] Create a new event")
		fmt.Println("[2] Register to an event as a volunteer")
		fmt.Println("[3] List all current events")
		fmt.Println("[4] List all jobs for a specific event")
		fmt.Println("[5] List the volunteers repartition for a specific event")
		fmt.Println("[6] To close an event")
		fmt.Println("[7] To terminate the process")

		choice = c.integerReader("")
		switch choice {
		case 1:
			createEvent(c)
		case 2:
			volunteerRegistration(c)
		case 3:
			printEvents(c)
		case 4:
			listJobs(c)
		case 5:
			volunteerRepartition(c)
		case 6:
			closeEvent(c)
		case 7:
			c.Close()
			return
		default:
			fmt.Println("You have entered a bad request")
		}
		fmt.Println()
	}
}

// loginClient allow the user to log in
// If the login is invalid, the user has to try again
func loginClient(c *Connection) {
	for {
		username := stringReader("Enter your username : ")
		password := stringReader("Enter your password : ")

		if c.LoginClient(username, password) {
			break
		}
	}
}

// createEvent creates a new event made by an organizer
// The user has to log in and must be an organizer
func createEvent(c *Connection) {
	loginClient(c)

	eventName := stringReader("Enter the event name : ")
	fmt.Println("List all job's name followed by the number of volunteers needed\n" +
		"(tap STOP when ended) : ")

	var jobList []string
	jobList = append(jobList, eventName)
	var i = 0
	for {
		i++
		jobName := stringReader("Insert a name for Job " + strconv.Itoa(i) + ": ")
		if strings.Compare(jobName, "STOP") == 0 {
			break
		}
		nbVolunteers := c.integerReader("Number of volunteers needed for this job : ")
		jobList = append(jobList, jobName, fmt.Sprint(nbVolunteers))
	}
	c.CreateEvent(jobList)
}

// printEvents prints all the events
func printEvents(c *Connection) {
	c.PrintEvents()
}

// volunteerRegistration allows a volunteer to register to an event
// The user has to log in
func volunteerRegistration(c *Connection) {
	loginClient(c)

	var eventId int
	var jobId int
	input := stringReader("Enter [event id] [job id] : ")
	_, err := fmt.Sscan(input, &eventId, &jobId)
	if err != nil {
		log.Fatal(err)
	}
	c.VolunteerRegistration(eventId, jobId)
}

// listJobs prints all the jobs for a specific event
func listJobs(c *Connection) {
	eventId := c.integerReader("Enter event id : ")
	c.ListJobs(eventId)
}

// volunteerRepartition prints the volunteer repartition for a specific event
func volunteerRepartition(c *Connection) {
	eventId := c.integerReader("Enter event id : ")
	c.VolunteerRepartition(eventId)
}

// closeEvent closes an event
// The user has to log in and must be the organizer of the event
func closeEvent(c *Connection) {
	loginClient(c)
	eventId := c.integerReader("Enter event id: ")
	c.CloseEvent(eventId)
}

// stringReader reads a string from the console
func stringReader(optionalMessage string) string {
	fmt.Print(optionalMessage)

	message, err := consoleIn.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimRight(message, "\r\n")
}

// integerReader reads an integer from the console and returns it
func (c *Connection) integerReader(optionalMessage string) int {
	for {
		fmt.Print(optionalMessage)
		n, err := strconv.ParseInt(stringReader(""), 10, 32)
		if errors.Is(err, strconv.ErrSyntax) {
			fmt.Println()
			fmt.Println("Not a number !")
		} else if errors.Is(err, strconv.ErrRange) {
			fmt.Println()
			fmt.Println("Number too big !")
		} else if err != nil {
			log.Fatal(err)
		} else {
			return int(n)
		}
	}
}

// printDataPacket prints the content of a data packet
func printDataPacket(data protocol.DataPacket) {
	for i := 0; i < len(data.Data); i++ {
		fmt.Println(data.Data[i])
	}
}
