package models

import (
	"SDR_Labo1/src/communication/server/tests"
	"errors"
	"reflect"
	"testing"
)

func TestCreateJob(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.Events[0].CreateJob("Test", 3)
	if err != nil {
		t.Error(err)
	}
	got := db.Events
	want := Events{
		{
			ID:        0,
			Name:      "Festival de la musique",
			Organizer: "Sarah Croche",
			Jobs: []Job{
				{
					ID:       0,
					Name:     "Buvette",
					Required: 2,
					Volunteers: []string{
						"Alex Terrieur",
					},
				},
				{
					ID:         1,
					Name:       "Sécurité",
					Required:   3,
					Volunteers: []string{},
				},
				{
					ID:         2,
					Name:       "Test",
					Required:   3,
					Volunteers: []string{},
				},
			},
		},
		{
			ID:        1,
			Name:      "Fête de la science",
			Organizer: "Ondine Akeleur",
			Jobs: []Job{
				{
					ID:         0,
					Name:       "Buvette",
					Required:   2,
					Volunteers: []string{},
				},
				{
					ID:         1,
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
	if !errors.Is(err, ErrorJobExists) {
		t.Error("Error not raised")
	}
}

func TestCreateJobErrorIfNameIsEmpty(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t)).Events
	_, err := db[0].CreateJob("", 3)
	if !errors.Is(err, ErrorJobNameEmpty) {
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
		ID:       0,
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
