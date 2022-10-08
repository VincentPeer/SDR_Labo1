package ui

import (
	"bufio"
	"fmt"
)

func ClientInterface(readWriter *bufio.ReadWriter) {
	fmt.Println("Welcome!\n Login : \n Please enter your username : ")
	//username, _ := readWriter.ReadString('\n') // todo gérer le 2ème attribut de retour (err)
	fmt.Println("Login : \n Please enter your password : ")
	//password, _ := readWriter.ReadString('\n')
	readWriter.WriteString("essai de co\n")
	//readWriter.WriteString(username + ", " + password) // todo check err

	for {
		fmt.Println("what do you want to write?")
		message, _ := readWriter.ReadString('\n')
		if message == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
		_, err := readWriter.WriteString(message)
		if err != nil {
			return
		}
	}

}
