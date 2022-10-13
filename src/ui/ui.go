package ui

import (
	"bufio"
	"fmt"
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

func loginClient(reader *bufio.Reader, writer *bufio.Writer) bool {
	fmt.Println("Enter your username : ")
	username, readUsernameError := reader.ReadString('\n') // todo gérer le 2ème attribut de retour (err)

	fmt.Println("Enter your password : ")
	password, readPasswordError := reader.ReadString('\n')

	_, writeError := writer.WriteString("LOG," + username + "," + password) // todo check err
	writer.Flush()
	response, responseError := reader.ReadString('\n')
	if readUsernameError != nil || readPasswordError != nil ||
		responseError != nil || writeError != nil || response == "FALSE" {
		return false
	} else {
		return true
	}
}

func CreateEvent(reader *bufio.Reader, writer *bufio.Writer) bool {
	fmt.Println("Enter the event name : ")
	eventName, _ := reader.ReadString('\n')
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
		writer.WriteString("EVENT," + eventName + ",")

	}
	return true
}
