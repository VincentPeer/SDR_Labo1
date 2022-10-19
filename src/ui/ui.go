package ui

import (
	"SDR_Labo1/src/communication/protocol"
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var messagingProtocol = &protocol.TcpProtocol{}

const EOF = "\r\n"

// Main function that communicate with a client, from here he can go through each
// functionnality offered on this service
func UserInterface(consoleReader *bufio.Reader, serverReader *bufio.Reader, serverWriter *bufio.Writer) {
	stop := protocol.DataPacket{Type: protocol.STOP}
	fmt.Println("Welcome!")

	var choice int
	for {
		fmt.Println("Choose one of the following functionnality")
		fmt.Println("[1] Create a new event")
		fmt.Println("[2] Register to an event as a volunteer")
		fmt.Println("[3] See all current events")
		fmt.Println("[4] See the volunteers repartiton for a specific event")
		fmt.Println("[5] To terminate the process")

		choice = integerReader(consoleReader)
		switch choice {
		case 1:
			createEvent(consoleReader, serverReader, serverWriter)
		case 2:
			volunteerRegistration()
		case 3:
			printEvents()
		case 4:
			volunteerRepartition()
		case 5:
			writeToServer(serverWriter, stop)
			return
		default:
			fmt.Println("You have entered a bad request")
		}
	}
}

// Gestion du login client
// Le client presse enter après chaque entrée, et ne doit pas saisir de ',' dans ses données
func loginClient(consoleReader *bufio.Reader, serverReader *bufio.Reader, serverWriter *bufio.Writer) bool {
	username := stringReader(consoleReader, "Enter your username : ")
	password := stringReader(consoleReader, "Enter your password : ")

	result := protocol.DataPacket{Type: protocol.LOGIN, Data: []string{username, password}}

	// Envoi formulaire de login
	writeToServer(serverWriter, result)

	// Traitement de la réponse après vérification du login par le serveur
	response := readFromServer(serverReader)
	if response.Type == protocol.OK {
		fmt.Println("Welcome " + username + "!")
		return true
	} else {
		fmt.Println("You have entered an invalid username or password")
		return false
	}
}

// Create a new event, it  asks the client to login and then the name of the event with each job's information
func createEvent(consoleReader *bufio.Reader, serverReader *bufio.Reader, serverWriter *bufio.Writer) bool {
	for {
		if loginClient(consoleReader, serverReader, serverWriter) == true {
			break
		}
	}

	eventName := stringReader(consoleReader, "Enter the event name : ")
	fmt.Println("List all job's name followed by the number of volunteers needed\n" +
		"(tap STOP when ended) : ")

	var jobList []string
	jobList = append(jobList, eventName)
	var i = 0
	for {
		i++
		jobName := stringReader(consoleReader, "Insert a name for Job "+strconv.Itoa(i)+": ")
		if strings.Compare(jobName, "STOP") == 0 {
			break
		}

		fmt.Print("Number of volunteers needed : ")
		nbVolunteers, err := strconv.ParseInt(stringReader(consoleReader, ""), 10, 32)
		if err != nil {
			log.Fatal(err)
		}

		jobList = append(jobList, jobName, fmt.Sprint(nbVolunteers))
	}

	eventResult := protocol.DataPacket{Type: protocol.CREATE_EVENT, Data: jobList}
	writeToServer(serverWriter, eventResult)
	response := readFromServer(serverReader)
	fmt.Println("response from server : ", response)
	if response.Type == protocol.OK {
		fmt.Println("Event registrated, thank you!")
		return true
	} else {
		return false
	}
}

func volunteerRegistration() {

}

func printEvents() {

}

func volunteerRepartition() {

}

func stringReader(reader *bufio.Reader, optionalMessage string) string {
	fmt.Print(optionalMessage)

	message, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}
	// Supression des retours à la ligne, et formatage pour l'envoi au serveur
	return strings.TrimRight(message, EOF)
}

func readFromServer(reader *bufio.Reader) protocol.DataPacket {
	message, err := reader.ReadString(protocol.DELIMITER)
	if err != nil {
		log.Fatal(err)
	}
	data, e := messagingProtocol.Receive(message)
	if e != nil {
		log.Fatal(e)
	}
	return data
}

func integerReader(reader *bufio.Reader) int {
	var n int
	nbScanned, err := fmt.Fscan(reader, &n)
	if err != nil {
		log.Fatal(err)
	} else if nbScanned != 1 {
		log.Fatal("Expected one argument, actual : " + strconv.Itoa(nbScanned))
	}
	reader.ReadString('\n') // clean the buffer
	return n
}

func writeToServer(writer *bufio.Writer, data protocol.DataPacket) {
	m, e := messagingProtocol.ToSend(data)
	fmt.Println(m)
	if e != nil {
		log.Fatal(e)
	}
	_, writtingError := writer.WriteString(m)
	if writtingError != nil {
		log.Fatal(writtingError)
	}
	flushError := writer.Flush()
	if flushError != nil {
		log.Fatal(flushError)
	}
}
