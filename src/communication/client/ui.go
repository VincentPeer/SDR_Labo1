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

// Main function that communicate with a connection, from here he can go through each
// functionnality offered on this service
func userInterface(c *connection) {
	stop := protocol.DataPacket{Type: protocol.STOP}
	fmt.Println("Welcome!")

	var choice int
	for {
		fmt.Println("Choose one of the following functionnality")
		fmt.Println("[1] Create a new event")
		fmt.Println("[2] Register to an event as a volunteer")
		fmt.Println("[3] List all current events")
		fmt.Println("[4] List the volunteers repartiton for a specific event")
		fmt.Println("[5] To terminate the process")

		choice = c.integerReader()
		switch choice {
		case 1:
			c.createEvent()
		case 2:
			c.volunteerRegistration()
		case 3:
			c.printEvents()
		case 4:
			c.volunteerRepartition()
		case 5:
			c.writeToServer(stop)
			return
		default:
			fmt.Println("You have entered a bad request")
		}
	}
}

// Gestion du login connection
// Le connection presse enter après chaque entrée, et ne doit pas saisir de ',' dans ses données
func (c *connection) loginClient() bool {
	username := c.stringReader("Enter your username : ")
	password := c.stringReader("Enter your password : ")

	result := protocol.DataPacket{Type: protocol.LOGIN, Data: []string{username, password}}

	// Envoi formulaire de login
	c.writeToServer(result)

	// Traitement de la réponse après vérification du login par le serveur
	response := c.readFromServer()
	if response.Type == protocol.OK {
		fmt.Println("Welcome " + username + "!")
		return true
	} else {
		fmt.Println("You have entered an invalid username or password")
		return false
	}
}

// Create a new event, it  asks the connection to login and then the name of the event with each job's information
func (c *connection) createEvent() bool {
	for {
		if c.loginClient() {
			break
		}
	}

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

func (c *connection) volunteerRegistration() bool {
	for {
		if c.loginClient() {
			break
		}
	}

	fmt.Println()

	return true
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
