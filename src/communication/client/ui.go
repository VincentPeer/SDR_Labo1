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

// userInterface is the main function that communicate with the server,
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
			c.createEvent()
		case 2:
			c.volunteerRegistration()
		case 3:
			c.printEvents()
		case 4:
			c.listJobs()
		case 5:
			c.volunteerRepartition()
		case 6:
			c.closeEvent()
		case 7:
			c.writeToServer(protocol.DataPacket{Type: protocol.STOP})
			return
		default:
			fmt.Println("You have entered a bad request")
		}
		fmt.Println()
	}
}

// loginClient ask the user to enter his username and password and check if the login is correct
// the client is asked to enter his username and password until the login is correct
func (c *connection) loginClient() {
	for {
		username := c.stringReader("Enter your username : ")
		password := c.stringReader("Enter your password : ")

		// Send the login request to the server
		login := protocol.DataPacket{Type: protocol.LOGIN, Data: []string{username, password}}

		// Manage the login request with the server
		response, _ := c.serverRequest(login, "Welcome "+username+"!")
		if response {
			break
		}
	}
}

// createEvent creates a new event makde by an organizer
func (c *connection) createEvent() bool {
	c.loginClient()

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
	event := protocol.DataPacket{Type: protocol.CREATE_EVENT, Data: jobList}
	response, _ := c.serverRequest(event, "Event registrated, thank you!")
	return response
}

func (c *connection) printEvents() {
	eventFound, data := c.serverRequest(protocol.DataPacket{Type: protocol.GET_EVENTS}, "")

	if eventFound {
		for i := 0; i < len(data.Data); i++ {
			fmt.Println(data.Data[i])
		}
	}
}

func (c *connection) volunteerRegistration() {
	c.loginClient()

	var eventId int
	var jobId int

	eventId = c.integerReader("Enter the event id : ")
	jobId = c.integerReader("Enter the job id : ")
	registration := protocol.DataPacket{Type: protocol.EVENT_REG, Data: []string{strconv.Itoa(eventId), strconv.Itoa(jobId)}}
	c.serverRequest(registration, "Registration successful")
}

func (c *connection) listJobs() {
	var eventId int
	var jobId int
	input := c.stringReader("Enter [event id] [job id] : ")
	_, err := fmt.Sscan(input, &eventId, &jobId)
	if err != nil {
		log.Fatal(err)
	}
	request := protocol.DataPacket{Type: protocol.GET_EVENTS, Data: []string{strconv.Itoa(eventId), strconv.Itoa(jobId)}}
	response, data := c.serverRequest(request, "")

	if response {
		for i := 0; i < len(data.Data); i++ {
			fmt.Println(data.Data[i])
		}
	} else {
		return
	}
}

func (c *connection) volunteerRepartition() {
	var eventId int
	eventId = c.integerReader("Enter the event id : ")
	request := protocol.DataPacket{Type: protocol.GET_JOBS, Data: []string{strconv.Itoa(eventId)}}
	response, data := c.serverRequest(request, "")

	if response {
		for i := 0; i < len(data.Data); i++ {
			fmt.Println(data.Data[i])
		}
	}
}

func (c *connection) closeEvent() {
	eventId := c.integerReader("Enter the id of the events you want to close : ")
	closeEvent := protocol.DataPacket{Type: protocol.CLOSE_EVENT, Data: []string{strconv.Itoa(eventId)}}
	c.serverRequest(closeEvent, "Event closed successfully")
}

func (c *connection) stringReader(optionalMessage string) string {
	fmt.Print(optionalMessage)

	message, err := c.consoleIn.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimRight(message, EOF)
}

func (c *connection) readFromServer() protocol.DataPacket {
	message, err := c.serverIn.ReadString(protocol.DELIMITER)
	if err != nil {
		log.Fatal(err)
	}
	data, e := messagingProtocol.Receive(message)
	if e != nil {
		log.Fatal(e)
	}
	return data
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
	c.consoleIn.ReadString('\n') // clean the buffer
	return n
}

func (c *connection) writeToServer(data protocol.DataPacket) {
	m, e := messagingProtocol.ToSend(data)
	if e != nil {
		log.Fatal(e)
	}
	_, writtingError := c.serverOut.WriteString(m)
	if writtingError != nil {
		log.Fatal(writtingError)
	}
	flushError := c.serverOut.Flush()
	if flushError != nil {
		log.Fatal(flushError)
	}
}

func (c *connection) serverRequest(data protocol.DataPacket, onSuccessMessage string) (bool, protocol.DataPacket) {
	c.writeToServer(data)
	response := c.readFromServer()
	if response.Type == protocol.OK {
		// fmt.Println(onSuccessMessage) todo
		return true, response
	} else {
		fmt.Println(response.Data)
		return false, response
	}
}
