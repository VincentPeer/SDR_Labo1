package main

import (
	"reflect"
	"testing"
)

func TestCreateEvent(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	testDb, err := createEvent(db, "Test", "3")
	if err != nil {
		t.Error(err)
	}
	got := testDb
	want := []Event{
		{
			Name:      "Festival de la musique",
			Organizer: "3",
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
			Organizer: "4",
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
			Organizer: "3",
			Jobs:      []Job{},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestCreateEventErrorIfIdExists(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	_, err := createEvent(db, "Festival de la musique", "3")
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestCreateEventErrorIfNameIsEmpty(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	_, err := createEvent(db, "", "3")
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestCreateEventErrorIfOrganizerDoesntExist(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	_, err := createEvent(db, "Test", "5")
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestCreateEventErrorIfOrganizerIsEmpty(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	_, err := createEvent(db, "Test", "")
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestGetEvent(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	testDb, err := getEvent(db, "1")
	if err != nil {
		t.Error(err)
	}
	got := testDb
	want := Event{
		Name:      "Festival de la musique",
		Organizer: "3",
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
	db := loadConfig(getTestData(t)).Events
	_, err := getEvent(db, "3")
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestCreateJob(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	testDb, err := createJob(db[0], "Test", 3)
	if err != nil {
		t.Error(err)
	}
	got := testDb
	want := []Event{
		{
			Name:      "Festival de la musique",
			Organizer: "3",
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
			Organizer: "4",
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
	db := loadConfig(getTestData(t)).Events
	_, err := createJob(db[0], "Buvette", 3)
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestCreateJobErrorIfNameIsEmpty(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	_, err := createJob(db[0], "", 3)
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestCreateJobErrorIfRequiredIsNegative(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	_, err := createJob(db[0], "Test", -1)
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestGetJob(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	testDb, err := getJob(db[0], "Buvette")
	if err != nil {
		t.Error(err)
	}
	got := testDb
	want := Job{
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
	_, err := getJob(db[0], "3")
	if err == nil {
		t.Error("Error not raised")
	}
}

func TestAddVolunteer(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	testDb, err := addVolunteer(db[0].Jobs[0], "2")
	if err != nil {
		t.Error(err)
	}
	got := testDb
	want := []Event{
		{
			Name:      "Festival de la musique",
			Organizer: "3",
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
			Organizer: "4",
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
	db := loadConfig(getTestData(t)).Events
	_, err := addVolunteer(db[0].Jobs[0], "1")
	if err == nil {
		t.Error("Error not raised")
	}
}
