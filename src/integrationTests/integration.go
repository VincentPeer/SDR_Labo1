package main

const (
	CONN_HOST        = "localhost"
	CONN_PORT        = "3333"
	CONFIG_FILE_PATH = "C:\\Users\\Vincent\\Documents\\Heig\\Semestre5\\SDR\\Laboratoires\\Labo01\\SDR_Labo1\\src\\integrationTests\\config.json"
)

/*


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
*/
