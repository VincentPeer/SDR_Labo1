package main

import (
	"SDR_Labo1/src/communication/server/models"
	"reflect"
	"testing"
)

func TestCreateEvent(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	testDb, err := db.CreateEvent("Test", "3")
	if err != nil {
		t.Error(err)
	}
	got := testDb
	want := []models.Event{
		{
			Name:      "Festival de la musique",
			Organizer: "3",
			Jobs: []models.Job{
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
			Organizer: "4",
			Jobs: []models.Job{
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
			Organizer: "3",
			Jobs:      []models.Job{},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestCreateEventErrorIfIdExists(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	_, err := db.CreateEvent("Festival de la musique", "3")
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestCreateEventErrorIfNameIsEmpty(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	_, err := db.CreateEvent("", "3")
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestCreateEventErrorIfOrganizerDoesntExist(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	_, err := db.CreateEvent("Test", "5")
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestCreateEventErrorIfOrganizerIsEmpty(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	_, err := db.CreateEvent("Test", "")
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestGetEvent(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	testDb, err := db.GetEvent("1")
	if err != nil {
		t.Error(err)
	}
	got := testDb
	want := models.Event{
		Name:      "Festival de la musique",
		Organizer: "3",
		Jobs: []models.Job{
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
	db := loadConfig(getTestData(t)).Events
	_, err := db.GetEvent("3")
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestCreateJob(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	testDb, err := db[0].CreateJob("Test", 3)
	if err != nil {
		t.Error(err)
	}
	got := testDb
	want := []models.Event{
		{
			Name:      "Festival de la musique",
			Organizer: "3",
			Jobs: []models.Job{
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
			Organizer: "4",
			Jobs: []models.Job{
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
	db := loadConfig(getTestData(t)).Events
	_, err := db[0].CreateJob("Buvette", 3)
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestCreateJobErrorIfNameIsEmpty(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	_, err := db[0].CreateJob("", 3)
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestCreateJobErrorIfRequiredIsNegative(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	_, err := db[0].CreateJob("Test", -1)
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestGetJob(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	testDb, err := db[0].GetJob("Buvette")
	if err != nil {
		t.Error(err)
	}
	got := testDb
	want := models.Job{
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
	db := loadConfig(getTestData(t)).Events
	_, err := db[0].GetJob("3")
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestAddVolunteer(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	testDb, err := db[0].Jobs[0].AddVolunteer("2")
	if err != nil {
		t.Error(err)
	}
	got := testDb
	want := []models.Event{
		{
			Name:      "Festival de la musique",
			Organizer: "3",
			Jobs: []models.Job{
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
			Organizer: "4",
			Jobs: []models.Job{
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
	db := loadConfig(getTestData(t)).Events
	_, err := db[0].Jobs[0].AddVolunteer("1")
	if err == nil {
		t.Error("Error not raised")
	}
}
