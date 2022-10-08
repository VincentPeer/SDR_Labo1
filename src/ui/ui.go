package ui

import (
	"bufio"
	"fmt"
)

func UserInterface(readWriter *bufio.ReadWriter) {
	fmt.Println("Welcome!\n Login : \nPlease enter your username : ")
	username, _ := readWriter.ReadString('\n') // todo gérer le 2ème attribut de retour (err)
	fmt.Println("Please enter your password : ")
	password, _ := readWriter.ReadString('\n')
	fmt.Println("username " + username + " password : " + password)
	readWriter.WriteString(username + ", " + password) // todo check err
	readWriter.Flush()

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
