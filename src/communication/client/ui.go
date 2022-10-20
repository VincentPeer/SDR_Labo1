package client

import (
	"SDR_Labo1/src/communication/protocol"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var messagingProtocol = &protocol.TcpProtocol{}

const EOF = "\r\n"

// userInterface is the function that communicate with the a,
// the client can go through each different functionnality
func userInterface(c *connection) {
	fmt.Println("Welcome!")

	var choice int
	for {
		fmt.Println("Choose one of the following functionnality")
		fmt.Println("[1] Create a new event")
		fmt.Println("[2] Register to an event as a volunteer")
		fmt.Println("[3] List all current events")
		fmt.Println("[4] List all jobs for a specific event")
		fmt.Println("[5] List the volunteers repartiton for a specific event")
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
			c.writeToServer(protocol.DataPacket{Type: protocol.STOP})
			return
		default:
			fmt.Println("You have entered a bad request")
		}
		fmt.Println()
	}
}

func loginClient(c *connection) {
	for {
		username := c.stringReader("Enter your username : ")
		password := c.stringReader("Enter your password : ")

		if c.LoginClient(username, password) {
			break
		}
	}
}

// createEvent creates a new event makde by an organizer
func createEvent(c *connection) bool {
	loginClient(c)

	eventName := c.stringReader("Enter the event name : ")
	fmt.Println("List all job's name followed by the number of volunteers needed\n" +
		"(tap STOP when ended) : ")

	var jobList []string
	jobList = append(jobList, eventName)
	var i = 0
	for {
		i++
		jobName := c.stringReader("Insert a name for Job " + strconv.Itoa(i) + ": ")
		if strings.Compare(jobName, "STOP") == 0 {
			break
		}

		fmt.Print("Number of volunteers needed : ")
		nbVolunteers, err := strconv.ParseInt(c.stringReader(""), 10, 32)
		if err != nil {
			log.Fatal(err)
		}

		jobList = append(jobList, jobName, fmt.Sprint(nbVolunteers))
	}
	return c.createEvent(jobList)
}

func printEvents(c *connection) {
	c.printEvents()
}

func volunteerRegistration(c *connection) {
	loginClient(c)

	var eventId int
	var jobId int
	input := c.stringReader("Enter [event id] [job id] : ")
	_, err := fmt.Sscan(input, &eventId, &jobId)
	if err != nil {
		log.Fatal(err)
	}
	c.volunteerRegistration(eventId, jobId)
}

func listJobs(c *connection) {
	eventId := c.integerReader("Enter event id : ")
	c.listJobs(eventId)
}

func volunteerRepartition(c *connection) {
	var eventId int
	eventId = c.integerReader("Enter event id : ")
	c.volunteerRepartition(eventId)
}

func closeEvent(c *connection) {
	loginClient(c)
	eventId := c.integerReader("Enter event id: ")
	c.closeEvent(eventId)
}

func (c *connection) stringReader(optionalMessage string) string {
	fmt.Print(optionalMessage)

	message, err := c.consoleIn.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimRight(message, EOF)
}

func readFromServer(c *connection) protocol.DataPacket {
	return c.readFromServer()
}

func (c *connection) integerReader(optionalMessage string) int {
	fmt.Print(optionalMessage)
	var n int
	nbScanned, err := fmt.Fscan(c.consoleIn, &n)
	if err != nil {
		log.Fatal(err)
	} else if nbScanned != 1 {
		log.Fatal("Expected one argument, actual : " + strconv.Itoa(nbScanned))
	}
	_, e := c.consoleIn.ReadString('\n') // clean the buffer
	if e != nil {
		log.Fatal(e)
	}
	return n
}

//func writeToServer(c *connection, data protocol.DataPacket) {
//	c.writeToServer(data)
//}

// serverRequest send a DataPacket to the server and return a boolean to know if the request was successful
// with a DataPacket containing the data response
//func serverRequest(c *connection, data protocol.DataPacket) (bool, protocol.DataPacket) {
//	return c.serverRequest(data)
//}

func printDataPacket(data protocol.DataPacket) {
	for i := 0; i < len(data.Data); i++ {
		fmt.Println(data.Data[i])
	}
}
