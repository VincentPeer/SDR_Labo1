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
	// loginTest(conn)

	// Create event tests
	//createEvent(conn)

	// List events tests
	//listEvents(conn)

	// List jobs tests
	// listJobs(conn)

	// List volunteer repartition  tests

	conn.Close()
}

func loginTest(conn *client.Connection) {
	fmt.Println("Test type : Login should fail with unknown user")
	loginStat := conn.LoginClient("unknown", "unknown")
	fmt.Println(loginStat, "\n")

	fmt.Println("Test type : Login should fail with wrong password")
	loginStat = conn.LoginClient("John", "wrong")
	fmt.Println(loginStat, "\n")

	fmt.Println("Test type : Login should fail with empty fields")
	loginStat = conn.LoginClient("", "")
	fmt.Println(loginStat, "\n")

	fmt.Println("Test type : Login should fail with wrong username")
	loginStat = conn.LoginClient("wrong", "123")
	fmt.Println(loginStat, "\n")

	fmt.Println("Test type : Login should succeed for a volunteer")
	loginStat = conn.LoginClient("James", "12345")
	fmt.Println(loginStat, "\n")

	fmt.Println("Test type : Login should succeed for an organizer")
	loginStat = conn.LoginClient("John", "123")
	fmt.Println(loginStat, "\n")
}

func createEvent(conn *client.Connection) {
	fmt.Println("Test type : Creation should fail without login")
	event := []string{"eventName", "job1", "2", "job2", "2"}
	loginStat := conn.CreateEvent(event)
	fmt.Println(loginStat, "\n")

	fmt.Println("Test type : Creation should fail with volunteer login")
	conn.LoginClient("James", "12345")
	event = []string{"eventName", "job1", "2", "job2", "2"}
	loginStat = conn.CreateEvent(event)
	fmt.Println(loginStat, "\n")

	fmt.Println("Test type : Creation should succeed with organizer login")
	conn.LoginClient("John", "123")
	event = []string{"eventName", "job1", "2", "job2", "2"}
	loginStat = conn.CreateEvent(event)
	fmt.Println(loginStat, "\n")

	fmt.Println("Test type : Creation should fail with empty fields")
	conn.LoginClient("John", "123")
	event = []string{""}
	loginStat = conn.CreateEvent(event)
	fmt.Println(loginStat, "\n")

	fmt.Println("Test type : Creation should succeed with name but without jobs")
	conn.LoginClient("John", "123")
	event = []string{"eventName"}
	loginStat = conn.CreateEvent(event)
	fmt.Println(loginStat, "\n")

	fmt.Println("Test type : Creation should succeed with name and one job")
	conn.LoginClient("John", "123")
	event = []string{"eventName", "job1", "3"}
	loginStat = conn.CreateEvent(event)
	fmt.Println(loginStat, "\n")

	fmt.Println("Test type : Creation should succeed with name and many jobs")
	conn.LoginClient("John", "123")
	event = []string{"eventName", "job1", "3", "job2", "2", "job3", "13"}
	loginStat = conn.CreateEvent(event)
	fmt.Println(loginStat, "\n")
}

func listEvents(conn *client.Connection) {
	fmt.Println("Test type : Current events list should succeed without login")
	conn.PrintEvents()
	fmt.Println("\n")

	fmt.Println("Test type : Current events list should succeed with login")
	conn.LoginClient("John", "123")
	conn.PrintEvents()
	fmt.Println("\n")

	fmt.Println("Test type : Current events list should add new event after its creation")
	event := []string{"AquaPoney", "buvette", "3", "entry", "2", "cleaning", "13"}
	conn.CreateEvent(event)
	conn.PrintEvents()
	fmt.Println("\n")

	fmt.Println("Test type : Current events list should mark as closed a close event")
	conn.LoginClient("Sarah Croche", "123456")
	conn.CloseEvent(0)
	conn.PrintEvents()
	fmt.Println("\n")
}

func listJobs(conn *client.Connection) {
	fmt.Println("Test type : Jobs list should succeed without login")
	conn.ListJobs(0)
	fmt.Println("\n")

	fmt.Println("Test type : Jobs list should succeed with login")
	conn.LoginClient("John", "123")
	conn.ListJobs(0)
	fmt.Println("\n")
}
