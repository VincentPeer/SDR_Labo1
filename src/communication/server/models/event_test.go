package models

import (
	"SDR_Labo1/src/communication/server/tests"
	"reflect"
	"testing"
)

func TestCreateEvent(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	testDb, err := db.CreateEvent("Test", "Sarah Croche")
	if err != nil {
		t.Error(err)
	}
	got := testDb.Events
	want := Events{
		{
			Name:      "Festival de la musique",
			Organizer: "Sarah Croche",
			Jobs: []Job{
				{
					Name:     "Buvette",
					Required: 2,
					Volunteers: []string{
						"Alex Terrieur",
					},
				},
				{
					Name:       "Sécurité",
					Required:   3,
					Volunteers: []string{},
				},
			},
		},
		{
			Name:      "Fête de la science",
			Organizer: "Ondine Akeleur",
			Jobs: []Job{
				{
					Name:       "Buvette",
					Required:   2,
					Volunteers: []string{},
				},
				{
					Name:       "Sécurité",
					Required:   3,
					Volunteers: []string{},
				},
			},
		},
		{
			Name:      "Test",
			Organizer: "Sarah Croche",
			Jobs:      []Job{},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestCreateEventErrorIfIdExists(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.CreateEvent("Festival de la musique", "Sarah Croche")
	if err != ErrorEventExists {
		t.Error("Error not raised")
	}
}

func TestCreateEventErrorIfNameIsEmpty(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.CreateEvent("", "Sarah Croche")
	if err != ErrorEventNameEmpty {
		t.Error("Error not raised")
	}
}

func TestCreateEventErrorIfOrganizerDoesntExist(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.CreateEvent("Test", "Bob")
	if err != ErrorUserNotFound {
		t.Error("Error not raised")
	}
}

func TestCreateEventErrorIfOrganizerIsEmpty(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.CreateEvent("Test", "")
	if err != ErrorOrganizerEmpty {
		t.Error("Error not raised")
	}
}

func TestGetEvent(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	testDb, err := db.GetEvent("Festival de la musique")
	if err != nil {
		t.Error(err)
	}
	got := testDb
	want := Event{
		Name:      "Festival de la musique",
		Organizer: "Sarah Croche",
		Jobs: []Job{
			{
				Name:     "Buvette",
				Required: 2,
				Volunteers: []string{
					"Alex Terrieur",
				},
			},
			{
				Name:       "Sécurité",
				Required:   3,
				Volunteers: []string{},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetEventErrorIfIdDoesntExist(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.GetEvent("Bob's party")
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestCreateJob(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.Events[0].CreateJob("Test", 3)
	if err != nil {
		t.Error(err)
	}
	got := db.Events
	want := Events{
		{
			Name:      "Festival de la musique",
			Organizer: "Sarah Croche",
			Jobs: []Job{
				{
					Name:     "Buvette",
					Required: 2,
					Volunteers: []string{
						"Alex Terrieur",
					},
				},
				{
					Name:       "Sécurité",
					Required:   3,
					Volunteers: []string{},
				},
				{
					Name:       "Test",
					Required:   3,
					Volunteers: []string{},
				},
			},
		},
		{
			Name:      "Fête de la science",
			Organizer: "Ondine Akeleur",
			Jobs: []Job{
				{
					Name:       "Buvette",
					Required:   2,
					Volunteers: []string{},
				},
				{
					Name:       "Sécurité",
					Required:   3,
					Volunteers: []string{},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestCreateJobErrorIfIdExists(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t)).Events
	_, err := db[0].CreateJob("Buvette", 3)
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestCreateJobErrorIfNameIsEmpty(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t)).Events
	_, err := db[0].CreateJob("", 3)
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestGetJob(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t)).Events
	testDb, err := db[0].GetJob("Buvette")
	if err != nil {
		t.Error(err)
	}
	got := testDb
	want := &Job{
		Name:     "Buvette",
		Required: 2,
		Volunteers: []string{
			"Alex Terrieur",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetJobErrorIfIdDoesntExist(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t)).Events
	_, err := db[0].GetJob("Test")
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestAddVolunteer(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.Events[0].Jobs[0].AddVolunteer("Alain Terrieur")
	if err != nil {
		t.Error(err)
	}
	got := db.Events
	want := Events{
		{
			Name:      "Festival de la musique",
			Organizer: "Sarah Croche",
			Jobs: []Job{
				{
					Name:     "Buvette",
					Required: 2,
					Volunteers: []string{
						"Alex Terrieur",
						"Alain Terrieur",
					},
				},
				{
					Name:       "Sécurité",
					Required:   3,
					Volunteers: []string{},
				},
			},
		},
		{
			Name:      "Fête de la science",
			Organizer: "Ondine Akeleur",
			Jobs: []Job{
				{
					Name:       "Buvette",
					Required:   2,
					Volunteers: []string{},
				},
				{
					Name:       "Sécurité",
					Required:   3,
					Volunteers: []string{},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestAddVolunteerErrorIfVolunteerAlreadyExists(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t)).Events
	_, err := db[0].Jobs[0].AddVolunteer("Alex Terrieur")
	if err != ErrorVolunteerExists {
		t.Error("Error not raised")
	}
}
