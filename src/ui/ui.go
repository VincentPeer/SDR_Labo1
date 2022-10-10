package ui

import (
	"bufio"
	"fmt"
)

func UserInterface(readWriter *bufio.ReadWriter) {
	fmt.Println("Welcome!")
	loginClient(readWriter)
	for {
		fmt.Println("what do you want to write?")
		message, _ := readWriter.ReadString('\n')
		if message == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
		_, err := readWriter.WriteString(message)
		readWriter.Flush()
		if err != nil {
			return
		}
	}

}

func loginClient(readWriter *bufio.ReadWriter) bool {
	fmt.Println("Please enter your username : ")
	username, readUsernameError := readWriter.ReadString('\n') // todo gérer le 2ème attribut de retour (err)

	fmt.Println("Please enter your password : ")
	password, readPasswordError := readWriter.ReadString('\n')

	//fmt.Println("username " + username + " password : " + password)

	_, writeError := readWriter.WriteString(username + ", " + password) // todo check err
	readWriter.Flush()
	if readUsernameError != nil || readPasswordError != nil || writeError != nil {
		return false
	} else {
		return true
	}
}
