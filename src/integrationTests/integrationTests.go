package main

import (
	"SDR_Labo1/src/communication/client"
	"SDR_Labo1/src/communication/server"
	"fmt"
)

func main() {
	fmt.Println("Running integration tests")

	go server.NewServer(false).Start()
	conn := client.CreateConnection(false)
	fmt.Println("Client created...")
	fmt.Println("Starting tests...")

	// Login tests
	loginTest(conn)

	// Registration tests
	volunteerRegistration(conn)

	// Create event tests
	createEvent(conn)

	// List events tests
	listEvents(conn)

	// List jobs tests
	listJobs(conn)

	// List volunteer repartition  tests
	volunteerRepartition(conn)

	// Close event tests
	closeEvent(conn)

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

func volunteerRegistration(conn *client.Connection) {
	// Registration tests with event 0, volunteer James
	fmt.Println("Test type : A volunteer should be able to register to an event")
	conn.LoginClient("James", "12345")
	conn.VolunteerRegistration(0, 0)
	conn.VolunteerRepartition(0)
	fmt.Println("\n")

	fmt.Println("Test type : Registering twice should remove the first registration and add to the second")
	conn.LoginClient("James", "12345")
	conn.VolunteerRegistration(0, 1)
	conn.VolunteerRepartition(0)
	fmt.Println("\n")

	fmt.Println("Test type : Registering should fail when the job is already full")
	conn.LoginClient("James", "12345")
	conn.VolunteerRegistration(0, 1)
	conn.VolunteerRepartition(0)
	fmt.Println("\n")

	fmt.Println("Test type : Registering should fail when the event is closed")
	conn.LoginClient("John", "123")
	conn.CloseEvent(3)
	conn.LoginClient("James", "12345")
	conn.VolunteerRegistration(3, 0)
	conn.VolunteerRepartition(3)
	fmt.Println("\n")

}

func createEvent(conn *client.Connection) {
	fmt.Println("Test type : Event creation should fail without login")
	event := []string{"eventName", "job1", "2", "job2", "2"}
	conn.CreateEvent(event)
	fmt.Println("\n")

	fmt.Println("Test type : Event creation should fail with volunteer login")
	conn.LoginClient("James", "12345")
	event = []string{"eventName", "job1", "2", "job2", "2"}
	conn.CreateEvent(event)
	fmt.Println("\n")

	fmt.Println("Test type : Event creation should succeed with organizer login")
	conn.LoginClient("John", "123")
	event = []string{"eventName", "job1", "2", "job2", "2"}
	conn.CreateEvent(event)
	fmt.Println("\n")

	fmt.Println("Test type : Event creation should fail with empty fields")
	conn.LoginClient("John", "123")
	event = []string{""}
	conn.CreateEvent(event)
	fmt.Println("\n")

	fmt.Println("Test type : Event creation should succeed with name but without jobs")
	conn.LoginClient("John", "123")
	event = []string{"eventName"}
	conn.CreateEvent(event)
	fmt.Println("\n")

	fmt.Println("Test type : Event creation should succeed with name and one job")
	conn.LoginClient("John", "123")
	event = []string{"eventName", "job1", "3"}
	conn.CreateEvent(event)
	fmt.Println("\n")

	fmt.Println("Test type : Event creation should succeed with name and many jobs")
	conn.LoginClient("John", "123")
	event = []string{"eventName", "job1", "3", "job2", "2", "job3", "13"}
	conn.CreateEvent(event)
	fmt.Println("\n")
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
	fmt.Println("Test type : Jobs list should fail with unexisting event")
	conn.ListJobs(100)
	fmt.Println("\n")

	fmt.Println("Test type : Jobs list should succeed without login")
	conn.ListJobs(0)
	fmt.Println("\n")

	fmt.Println("Test type : Jobs list should succeed with login")
	conn.LoginClient("John", "123")
	conn.ListJobs(0)
	fmt.Println("\n")
}

func volunteerRepartition(conn *client.Connection) {
	fmt.Println("Test type : Volunteer repartition should fail with invalid event id")
	conn.VolunteerRepartition(100)
	fmt.Println("\n")

	fmt.Println("Test type : Volunteer repartition should print repartition for a given event")
	conn.VolunteerRepartition(0)
	fmt.Println("\n")
}

func closeEvent(conn *client.Connection) {
	fmt.Println("Test type : Closing an event should fail if the user isn't logged in")
	conn.CloseEvent(0)
	fmt.Println("\n")

	fmt.Println("Test type : Closing an event should fail if the user isn't the organizer")
	conn.LoginClient("James", "12345")
	conn.CloseEvent(0)
	fmt.Println("\n")

	fmt.Println("Test type : Closing an event should fail when closing unexist event")
	fmt.Println("Closing event 55")
	conn.LoginClient("Sarah Croche", "123456")
	conn.CloseEvent(55)
	fmt.Println("\n")

	fmt.Println("Test type : Closing an event should succeed if the user is the organizer")
	fmt.Println("Closing event 5")
	conn.PrintEvents()
	conn.LoginClient("Alain Du Bois", "123456")
	conn.CloseEvent(2)
	conn.PrintEvents()
	fmt.Println("\n")

	fmt.Println("Test type : Closing an event should fail when closing unexist event")
	fmt.Println("Closing event 55")
	conn.LoginClient("Sarah Croche", "123456")
	conn.CloseEvent(55)
	fmt.Println("\n")
}
