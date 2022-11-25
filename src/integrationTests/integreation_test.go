package integration_test

import (
	"SDR_Labo1/src/communication/client"
	"testing"
)

var (
	connHost = "localhost"
	connPort = "3333"
	conn     = client.CreateConnection(connHost, connPort, false)
)

func TestLoginClient(t *testing.T) {
	// Should fail with unknown user
	got := conn.LoginClient("unknown", "unknown")
	if got {
		t.Errorf("LoginClient() Should fail with unknown user = %v, want %v", got, false)
	}

	// Should fail with no password
	got = conn.LoginClient("unknown", "")
	if got {
		t.Errorf("LoginClient() Should fail with no password = %v, want %v", got, false)
	}

	// Should fail with no username
	got = conn.LoginClient("", "unknown")
	if got {
		t.Errorf("LoginClient() Should fail with no username = %v, want %v", got, false)
	}

	// Should fail with empty fields
	got = conn.LoginClient("", "")
	if got {
		t.Errorf("LoginClient() Should fail with empty fields = %v, want %v", got, false)
	}

	// Should fail with a wrong password
	got = conn.LoginClient("John", "wrong")
	if got {
		t.Errorf("LoginClient() Should fail with a wrong password = %v, want %v", got, false)
	}

	// Should succeed for an organizer
	got = conn.LoginClient("John", "123")
	if !got {
		t.Errorf("LoginClient() Should succeed for an organizer = %v, want %v", got, true)
	}

	// Should succeed for a volunteer
	got = conn.LoginClient("James", "12345")
	if !got {
		t.Errorf("LoginClient() Should succeed for a volunteer = %v, want %v", got, true)
	}
}

func TestCreateEvent(t *testing.T) {
	// Should fail with no login
	event := []string{"eventName", "job1", "2", "job2", "2"}
	got := conn.CreateEvent(event)
	if got {
		t.Errorf("CreateEvent() Should fail with no login = %v, want %v", got, false)
	}

	// Should fail as volunteer logged
	conn.LoginClient("James", "12345") // Login as a volunteer
	got = conn.CreateEvent(event)
	if got {
		t.Errorf("CreateEvent() Should fail as volunteer logged = %v, want %v", got, false)
	}

	// Should succeed with valid data, logged as an organizer
	conn.LoginClient("John", "123") // Login as an organizer
	event = []string{"eventName", "job1", "2", "job2", "2"}
	got = conn.CreateEvent(event)
	if !got {
		t.Errorf("CreateEvent() Should succeed with valid data, logged as an organizer = %v, want %v", got, true)
	}

	// Should succeed without jobs (but with an event name)
	conn.LoginClient("John", "123") // Login as an organizer
	event = []string{"1"}
	got = conn.CreateEvent(event)
	if !got {
		t.Errorf("CreateEvent() Should succeed without jobs = %v, want %v", got, true)
	}
}

func TestListEvents(t *testing.T) {
	// Should succeed without login
	got := conn.PrintEvents()
	if !got {
		t.Errorf("ListEvents() Should succeed without login = %v, want %v", got, true)
	}

	// Should succeed being logged
	conn.LoginClient("John", "123")
	got = conn.PrintEvents()
	if !got {
		t.Errorf("ListEvents() Should succeed being logged = %v, want %v", got, true)
	}
}

func TestListJobs(t *testing.T) {
	// Should succeed without being logged
	got := conn.ListJobs(0)
	if !got {
		t.Errorf("ListJobs() Should succeed without being logged = %v, want %v", got, true)
	}

	// Should succeed being logged
	conn.LoginClient("John", "123")
	got = conn.ListJobs(0)
	if !got {
		t.Errorf("ListJobs() = %v, want %v", got, true)
	}

	// Should fail with unexisting event
	got = conn.ListJobs(1000)
	if got {
		t.Errorf("ListJobs() = %v, want %v", got, false)
	}
}

func TestVolunteerRegistration(t *testing.T) {
	// Registration should succeed with a volunteer registration to an existed event
	conn.LoginClient("James", "12345")
	got := conn.VolunteerRegistration(0, 0)
	if !got {
		t.Errorf("VolunteerRegistration() should succeed with a volunteer registration = %v, want %v", got, true)
	}

	// Registration should fail when the job is already full
	conn.LoginClient("James", "12345")
	got = conn.VolunteerRegistration(2, 1)
	if got {
		t.Errorf("VolunteerRegistration() should fail when the job is already full = %v, want %v", got, false)
	}

	// Registration should fail when the event is closed
	conn.LoginClient("James", "12345")
	got = conn.VolunteerRegistration(3, 0)
	if got {
		t.Errorf("VolunteerRegistration() should fail when the event is closed = %v, want %v", got, false)
	}

	// An organizer registration should fail
	conn.LoginClient("John", "123")
	got = conn.VolunteerRegistration(3, 1)
	if got {
		t.Errorf("VolunteerRegistration() organizer registration should fail = %v, want %v", got, false)
	}
}

func TestVolunteerRepartition(t *testing.T) {
	// Volunteer repartition should fail with invalid event id
	got := conn.VolunteerRepartition(100)
	if got {
		t.Errorf("VolunteerRepartition() repartition should fail with invalid event id = %v, want %v", got, false)
	}

	// Volunteer repartition should succeed for a valid id
	got = conn.VolunteerRepartition(1)
	if !got {
		t.Errorf("VolunteerRepartition() should succeed for a valid id = %v, want %v", got, true)
	}
}

func TestCloseEvent(t *testing.T) {
	// Closing an event should fail if the user isn't logged
	got := conn.CloseEvent(0)
	if got {
		t.Errorf("CloseEvent() should fail if the user isn't logged = %v, want %v", got, false)
	}

	// Closing an event should fail when closing non existed event
	conn.LoginClient("Sarah Croche", "123456")
	got = conn.CloseEvent(55)
	if got {
		t.Errorf("CloseEvent() should fail when closing non existed event = %v, want %v", got, false)
	}

	// Closing an event should fail if the user isn't the organizer
	conn.LoginClient("James", "12345")
	got = conn.CloseEvent(0)
	if got {
		t.Errorf("CloseEvent() should fail if the user isn't the organizer = %v, want %v", got, false)
	}

	// Closing an event should succeed if the user is the organizer
	conn.LoginClient("Alain Du Bois", "123456")
	got = conn.CloseEvent(2)
	if !got {
		t.Errorf("CloseEvent() should fail if the user isn't the organizer = %v, want %v", got, true)
	}

}
