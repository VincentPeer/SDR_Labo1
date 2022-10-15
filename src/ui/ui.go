package ui

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

func UserInterface(reader *bufio.Reader, writer *bufio.Writer) {
	fmt.Println("Welcome!")

	for {
		if !loginClient(reader, writer) {
			continue
		}

		message, _ := reader.ReadString('\n')
		if message == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}

// Gestion du login client
// Le client presse enter après chaque entrée, et ne doit pas saisir de ',' dans ses données
func loginClient(reader *bufio.Reader, writer *bufio.Writer) bool {
	username := readFromServer(reader, "Enter your username : ")
	password := readFromServer(reader, "Enter your password : ")

	// Supression des retours à la ligne, et formatage pour l'envoi au serveur
	username = strings.TrimSuffix(username, "\r\n") + ","
	password = strings.TrimSuffix(password, "\r\n")
	result := "LOGIN," + username + password + ";"

	// Envoi formulaire de login
	writeToServer(writer, result)

	// Traitement de la réponse après vérification du login par le serveur
	response := readFromServer(reader, "")
	fmt.Println("response from server is : " + response)
	if strings.Compare(response, "OK") == 0 {
		fmt.Println("Hello " + username + "!")
		return true
	} else {
		fmt.Println("You have entered an invalid username or password")
		return false
	}
}

func createEvent(reader *bufio.Reader, writer *bufio.Writer) bool {
	eventName := readFromServer(reader, "Enter the event name : ")
	var jobList []string
	fmt.Println("List all job's name followed by the number of volunteers needed\n" +
		"(tap double enter when ended) : ")
	for {
		jobName, jobError := reader.ReadString('\n')
		_, writeError := writer.WriteString("JOB," + jobName)
		if jobError != nil || writeError != nil {
			return false
		}
		jobList = append(jobList, jobName)
		m, _ := reader.ReadString('\n')
		if m == "\n" {
			break
		}
		writer.WriteString("CREATE_EVENT," + eventName + ",")
	}
	return true
}

func readFromServer(reader *bufio.Reader, textForUser string) string {
	fmt.Print(textForUser)
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
