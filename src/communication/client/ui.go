package client

import (
	"SDR_Labo1/src/communication/protocol"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var messagingProtocol = &protocol.TcpProtocol{}

const EOF = "\r\n"

var consoleOut = bufio.NewWriter(os.Stdin)
var consoleIn = bufio.NewReader(os.Stdin)

// UserInterface is the function that communicate with the a,
// the client can go through each different functionnality
func UserInterface(c *Connection) {
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
			c.Close()
			return
		default:
			fmt.Println("You have entered a bad request")
		}
		fmt.Println()
	}
}

func loginClient(c *Connection) {
	for {
		username := stringReader("Enter your username : ")
		password := stringReader("Enter your password : ")

		if c.LoginClient(username, password) {
			break
		}
	}
}

// createEvent creates a new event makde by an organizer
func createEvent(c *Connection) bool {
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

		fmt.Print("Number of volunteers needed : ")
		nbVolunteers, err := strconv.ParseInt(stringReader(""), 10, 32)
		if err != nil {
			log.Fatal(err)
		}

		jobList = append(jobList, jobName, fmt.Sprint(nbVolunteers))
	}
	return c.CreateEvent(jobList)
}

func printEvents(c *Connection) {
	c.PrintEvents()
}

func volunteerRegistration(c *Connection) {
	loginClient(c)

	var eventId int
	var jobId int
	input := stringReader("Enter [event id] [job id] : ")
	_, err := fmt.Sscan(input, &eventId, &jobId)
	if err != nil {
		log.Fatal(err)
	}
	c.volunteerRegistration(eventId, jobId)
}

func listJobs(c *Connection) {
	eventId := c.integerReader("Enter event id : ")
	c.ListJobs(eventId)
}

func volunteerRepartition(c *Connection) {
	var eventId int
	eventId = c.integerReader("Enter event id : ")
	c.VolunteerRepartition(eventId)
}

func closeEvent(c *Connection) {
	loginClient(c)
	eventId := c.integerReader("Enter event id: ")
	c.CloseEvent(eventId)
}

func stringReader(optionalMessage string) string {
	fmt.Print(optionalMessage)

	message, err := consoleIn.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimRight(message, EOF)
}

func readFromServer(c *Connection) protocol.DataPacket {
	return c.readFromServer()
}

func (c *Connection) integerReader(optionalMessage string) int {
	fmt.Print(optionalMessage)
	var n int
	nbScanned, err := fmt.Fscan(consoleIn, &n)
	if err != nil {
		log.Fatal(err)
	} else if nbScanned != 1 {
		log.Fatal("Expected one argument, actual : " + strconv.Itoa(nbScanned))
	}
	_, e := consoleIn.ReadString('\n') // clean the buffer
	if e != nil {
		log.Fatal(e)
	}
	return n
}

//func writeToServer(c *Connection, data protocol.DataPacket) {
//	c.writeToServer(data)
//}

// serverRequest send a DataPacket to the server and return a boolean to know if the request was successful
// with a DataPacket containing the data response
//func serverRequest(c *Connection, data protocol.DataPacket) (bool, protocol.DataPacket) {
//	return c.serverRequest(data)
//}

func printDataPacket(data protocol.DataPacket) {
	for i := 0; i < len(data.Data); i++ {
		fmt.Println(data.Data[i])
	}
}
