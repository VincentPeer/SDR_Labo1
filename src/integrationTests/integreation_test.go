package main

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
		t.Errorf("LoginClient() = %v, want %v", got, false)
	}

	// Should fail with no password
	got = conn.LoginClient("unknown", "")
	if got {
		t.Errorf("LoginClient() = %v, want %v", got, false)
	}

	// Should fail with no username
	got = conn.LoginClient("", "unknown")
	if got {
		t.Errorf("LoginClient() = %v, want %v", got, false)
	}

	// Should fail with empty fields
	got = conn.LoginClient("", "")
	if got {
		t.Errorf("LoginClient() = %v, want %v", got, false)
	}

	// Should fail with a wrong password
	got = conn.LoginClient("John", "wrong")
	if got {
		t.Errorf("LoginClient() = %v, want %v", got, false)
	}

	// Should succeed for an organizer
	got = conn.LoginClient("John", "123")
	if !got {
		t.Errorf("LoginClient() = %v, want %v", got, true)
	}

	// Should succeed for a volunteer
	got = conn.LoginClient("James", "12345")
	if !got {
		t.Errorf("LoginClient() = %v, want %v", got, true)
	}
}

func TestCreateEvent(t *testing.T) {
	// Should fail with no login
	event := []string{"eventName", "job1", "2", "job2", "2"}
	got := conn.CreateEvent(event)
	if got {
		t.Errorf("CreateEvent() = %v, want %v", got, false)
	}

	// Should fail as volunteer logged
	conn.LoginClient("James", "12345") // Login as a volunteer
	got = conn.CreateEvent(event)
	if got {
		t.Errorf("CreateEvent() = %v, want %v", got, false)
	}

	// Should succeed with valid data, logged as an organizer
	conn.LoginClient("John", "123") // Login as an organizer
	event = []string{"eventName", "job1", "2", "job2", "2"}
	got = conn.CreateEvent(event)
	if !got {
		t.Errorf("CreateEvent() = %v, want %v", got, true)
	}

	// Should succeed without jobs (but with an event name)
	conn.LoginClient("John", "123") // Login as an organizer
	event = []string{"1"}
	got = conn.CreateEvent(event)
	if !got {
		t.Errorf("CreateEvent() = %v, want %v", got, true)
	}
}

func TestListEvents(t *testing.T) {
	// Should succeed without login
	got := conn.PrintEvents()
	if !got {
		t.Errorf("CreateEvent() = %v, want %v", got, true)
	}

	// Should succeed being logged
	conn.LoginClient("John", "123")
	got = conn.PrintEvents()
	if !got {
		t.Errorf("CreateEvent() = %v, want %v", got, true)
	}
}

func TestListJobs(t *testing.T) {
	// Should succeed without being logged
	got := conn.ListJobs(0)
	if !got {
		t.Errorf("CreateEvent() = %v, want %v", got, true)
	}

	// Should succeed being logged
	conn.LoginClient("John", "123")
	got = conn.ListJobs(0)
	if !got {
		t.Errorf("CreateEvent() = %v, want %v", got, true)
	}

	// Should fail with unexisting event
	got = conn.ListJobs(1000)
	if got {
		t.Errorf("CreateEvent() = %v, want %v", got, false)
	}
}

func TestVolunteerRepartition(t *testing.T) {
	// A volunteer should be able to register to an event
	conn.LoginClient("James", "12345")
	conn.VolunteerRegistration(0, 0)
}
