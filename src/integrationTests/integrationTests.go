package main

import (
	"SDR_Labo1/src/communication/client"
	"SDR_Labo1/src/communication/server"
	"fmt"
)

func main() {
	fmt.Println("Running integration tests")

	go server.NewServer().Start()
	conn := client.CreateConnection()
	fmt.Println("Client created...")
	fmt.Println("Starting tests...")

	// Login tests
	loginTest(conn)

	conn.Close()
}

func loginTest(conn *client.Connection) {
	fmt.Println("Login fail with unknown user")
	loginStat := conn.LoginClient("unknown", "unknown")
	fmt.Println(loginStat)

	fmt.Println("Login fail with wrong password")
	loginStat = conn.LoginClient("John", "wrong")
	fmt.Println(loginStat)

	fmt.Println("Login fail with wrong username")
	loginStat = conn.LoginClient("wrong", "123")
	fmt.Println(loginStat)

	fmt.Println("Login success for a volunteer")
	loginStat = conn.LoginClient("James", "12345")
	fmt.Println(loginStat)

	fmt.Println("Login success for an organizer")
	loginStat = conn.LoginClient("John", "123")
	fmt.Println(loginStat)

}
