package ui

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const EOF = "\r\n"

func UserInterface(consoleReader *bufio.Reader, serverReader *bufio.Reader, serverWriter *bufio.Writer) {
	fmt.Println("Welcome!")

	for {
		fmt.Println("start loop")
		if loginClient(consoleReader, serverReader, serverWriter) == true {
			break
		}
	}
	createEvent(consoleReader, serverReader, serverWriter)

}

// Gestion du login client
// Le client presse enter après chaque entrée, et ne doit pas saisir de ',' dans ses données
func loginClient(consoleReader *bufio.Reader, serverReader *bufio.Reader, serverWriter *bufio.Writer) bool {
	username := stringReader(consoleReader, "Enter your username : ")
	password := stringReader(consoleReader, "Enter your password : ")

	// Supression des retours à la ligne, et formatage pour l'envoi au serveur
	username = strings.TrimSuffix(username, "\r\n") + ","
	password = strings.TrimSuffix(password, "\r\n")
	result := "LOGIN," + username + password + ";"

	// Envoi formulaire de login
	writeToServer(serverWriter, result)

	// Traitement de la réponse après vérification du login par le serveur
	response := stringReader(serverReader, "")
	fmt.Println("response from server is : " + response)
	if strings.Compare(response, "OK\n") == 0 {
		fmt.Println("Hello " + username + "!")
		return true
	} else {
		fmt.Println("You have entered an invalid username or password")
		return false
	}
}

func createEvent(consoleReader *bufio.Reader, serverReader *bufio.Reader, serverWriter *bufio.Writer) bool {
	eventName := stringReader(consoleReader, "Enter the event name : ")
	fmt.Println("List all job's name followed by the number of volunteers needed\n" +
		"(tap enter with empty field when ended) : ")

	var jobList []string
	var i = 0
	var nbVolunteers int
	for {
		i++
		jobName := stringReader(consoleReader, "Insert name for Job "+strconv.Itoa(i)+": ")
		if strings.Compare(jobName, EOF) == 0 {
			break
		}
		fmt.Print("Number of volunteers needed : ")
		fmt.Scanf("%d", &nbVolunteers)
		jobList = append(jobList, jobName+","+strconv.Itoa(nbVolunteers))
	}

	jobStringList := strings.Join(jobList, ", ")

	eventResult := "CREATE_EVENT," + eventName + "," + jobStringList + ";"
	eventResult = strings.Replace(eventResult, EOF, "", -1)
	fmt.Println(eventResult)
	writeToServer(serverWriter, eventResult)
	response := stringReader(serverReader, "")
	if strings.Compare(response, "OK"+"\n") == 0 {
		fmt.Println("Event registrated, thank you!")
		return true
	} else {
		return false
	}
}

func stringReader(reader *bufio.Reader, optinalMessage string) string {
	fmt.Print(optinalMessage)
	message, err := reader.ReadString('\n')
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
