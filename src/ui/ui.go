package ui

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func UserInterface(reader *bufio.Reader, writer *bufio.Writer, serverReader *bufio.Reader) {
	fmt.Println("Welcome!")

	for {
		if !loginClient(reader, writer, serverReader) {
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
func loginClient(reader *bufio.Reader, writer *bufio.Writer, serverReader *bufio.Reader) bool {
	username := readFromServer(reader, "Enter your username : ")
	password := readFromServer(reader, "Enter your password : ")

	// Supression des retours à la ligne, et formatage pour l'envoi au serveur
	username = strings.TrimSuffix(username, "\r\n") + ","
	password = strings.TrimSuffix(password, "\r\n")
	result := "LOGIN," + username + password + ";"

	// Envoi formulaire de login
	writeToServer(writer, result)

	// Traitement de la réponse après vérification du login par le serveur
	response := readFromServer(serverReader, "")
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
	var i = 0
	fmt.Println("List all job's name followed by the number of volunteers needed\n" +
		"(tap double enter when ended) : ")
	for {
		i++
		jobName := readFromServer(reader, "Insert name of Job "+strconv.Itoa(i)+": ")
		if jobName == "\n" {
			break
		}
		jobList = append(jobList, jobName)
	}
	writeToServer(writer, "CREATE_EVENT,"+eventName+";")
	response := readFromServer(reader, "")
	if strings.Compare(response, "OK") == 0 {
		return true
	} else {
		return false
	}
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
