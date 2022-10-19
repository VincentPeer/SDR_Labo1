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
		fmt.Println("[6] To terminate the process")

		choice = c.integerReader()
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
			c.writeToServer(protocol.DataPacket{Type: protocol.STOP})
			return
		default:
			fmt.Println("You have entered a bad request")
		}
	}
}

// loginClient ask the user to enter his username and password and check if the login is correct
// the client is asked to enter his username and password until the login is correct
func (c *connection) loginClient() {
	for {
		username := c.stringReader("Enter your username : ")
		password := c.stringReader("Enter your password : ")

		// Send the login request to the server
		result := protocol.DataPacket{Type: protocol.LOGIN, Data: []string{username, password}}
		c.writeToServer(result)

		// Read the response from the server and treat it
		response := c.readFromServer()
		if response.Type == protocol.OK {
			fmt.Println("Welcome " + username + "!")
			return
		} else {
			fmt.Println("You have entered an invalid username or password")
			continue
		}
	}
}

// createEvent creates a new event with an organizer
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

	eventResult := protocol.DataPacket{Type: protocol.CREATE_EVENT, Data: jobList}
	c.writeToServer(eventResult)
	response := c.readFromServer()
	fmt.Println("response from server : ", response)
	if response.Type == protocol.OK {
		fmt.Println("Event registrated, thank you!")
		return true
	} else {
		return false
	}
}

func (c *connection) printEvents() {
	c.writeToServer(protocol.DataPacket{Type: protocol.GET_EVENTS})
	response := c.readFromServer()
	if response.Type == protocol.OK {
		for i := 0; i < len(response.Data); i++ {
			fmt.Println(response.Data[i])
		}
	} else {
		fmt.Println("No event found")
	}
}

func (c *connection) volunteerRegistration() {
	c.loginClient()

	var eventId int
	var jobId int
	for {
		fmt.Println("Choose one of the following functionnality")
		fmt.Println("[1] List all open events with their jobs")
		fmt.Println("[2] Register to an event as a volunteer")
		fmt.Println("[3] To terminate the process")

		choice := c.integerReader()
		switch choice {
		case 1:
			c.printEvents()
		case 2:
			fmt.Println("Enter the id of the event you want to register to : ")
			eventId = c.integerReader()
			fmt.Println("Enter the id of the job you want to register to : ")
			jobId = c.integerReader()
			registration := protocol.DataPacket{Type: protocol.EVENT_REG, Data: []string{strconv.Itoa(eventId), strconv.Itoa(jobId)}}
			c.writeToServer(registration)
			response := c.readFromServer()
			if response.Type == protocol.OK {
				fmt.Println("Registration successful")
				return
			} else {
				fmt.Println("Registration failed")
				continue
			}
		case 3:
			c.writeToServer(protocol.DataPacket{Type: protocol.STOP})
			return
		}
	}
}

func (c *connection) listJobs() {
	var eventId int
	fmt.Println("Enter the jobsToList id : ")
	eventId = c.integerReader()
	jobsToList := protocol.DataPacket{Type: protocol.GET_JOBS, Data: []string{strconv.Itoa(eventId)}}
	c.writeToServer(jobsToList)
	response := c.readFromServer()
	if response.Type == protocol.OK {
		for i := 0; i < len(response.Data); i++ {
			fmt.Println(response.Data[i])
		}
	} else {
		fmt.Println("No job found")
	}
}

func (c *connection) volunteerRepartition() {

}

func (c *connection) stringReader(optionalMessage string) string {
	fmt.Print(optionalMessage)

	message, err := c.consoleIn.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}
	// Supression des retours à la ligne, et formatage pour l'envoi au serveur
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

func (c *connection) integerReader() int {
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
	fmt.Println(m)
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
