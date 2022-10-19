package ui

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
)

const EOF = "\r\n"

// Main function that communicate with a client, from here he can go through each
// functionnality offered on this service
func UserInterface(consoleReader *bufio.Reader, serverReader *bufio.Reader, serverWriter *bufio.Writer) {
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
			return // todo: ou break
		default:
			fmt.Println("You have entered a bad request")
		}
	}
}

// Gestion du login client
// Le client presse enter après chaque entrée, et ne doit pas saisir de ',' dans ses données
func loginClient(consoleReader *bufio.Reader, serverReader *bufio.Reader, serverWriter *bufio.Writer) bool {
	const LOG_PROTOCOL = "LOGIN"
	username := stringReader(consoleReader, "Enter your username : ")
	password := stringReader(consoleReader, "Enter your password : ")

	// Supression des retours à la ligne, et formatage pour l'envoi au serveur
	username = strings.TrimSuffix(username, EOF)
	password = strings.TrimSuffix(password, EOF)
	tmpResult := []string{LOG_PROTOCOL, username, password}
	result := strings.Join(tmpResult, ",") + ";"

	// Envoi formulaire de login
	writeToServer(serverWriter, result)

	// Traitement de la réponse après vérification du login par le serveur
	response := readFromServer(serverReader, "")
	if strings.Compare(response, "OK;") == 0 {
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
	var i = 0
	var nbVolunteers int
	for {
		i++
		jobName := stringReader(consoleReader, "Insert a name for Job "+strconv.Itoa(i)+": ")
		if strings.Compare(jobName, "STOP") == 0 {
			break
		}

		fmt.Print("Number of volunteers needed : ")
		fmt.Scanf("%d", &nbVolunteers)
		fmt.Println("read : ", nbVolunteers)

		jobName = strings.TrimSuffix(jobName, "\r\n")

		jobList = append(jobList, jobName+","+strconv.Itoa(nbVolunteers))
	}
	var jobStringList string
	if len(jobList) > 1 {
		jobStringList = strings.Join(jobList, ", ")
	} else if len(jobList) == 1 {
		jobStringList = jobList[0]
	}

	eventResult := "CREATE_EVENT," + eventName + "," + jobStringList + ";"
	eventResult = strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, eventResult)
	fmt.Println(eventResult)
	writeToServer(serverWriter, eventResult)
	response := readFromServer(serverReader, "")
	if strings.Compare(response, "OK"+";") == 0 {
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
	var s string

	nbScanned, err := fmt.Fscan(reader, &s)
	if err != nil {
		log.Fatal(err)
	} else if nbScanned != 1 {
		log.Fatal("Expected one argument, actual : " + strconv.Itoa(nbScanned))
	}
	reader.ReadString('\n') // clean the buffer
	return s
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

func readFromServer(reader *bufio.Reader, optinalMessage string) string {
	fmt.Print(optinalMessage)
	message, err := reader.ReadString(';')
	if err != nil {
		log.Fatal(err)
	}
	return message
}

func writeToServer(writer *bufio.Writer, message string) {
	_, writtingError := writer.WriteString(message)
	if writtingError != nil {
		log.Fatal(writtingError)
	}
	flushError := writer.Flush()
	if flushError != nil {
		log.Fatal(flushError)
	}
}
