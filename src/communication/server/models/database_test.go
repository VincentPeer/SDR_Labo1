package models

import (
	"SDR_Labo1/src/communication/server/tests"
	"errors"
	"reflect"
	"testing"
)

func TestLoadUsers(t *testing.T) {

	got := LoadDatabaseFromJson(tests.GetTestData(t)).Users
	want := map[string]User{
		"Alex Terrieur": {
			Name:     "Alex Terrieur",
			Password: "AlexPWD",
			Function: "volunteer",
		},
		"Alain Terrieur": {
			Name:     "Alain Terrieur",
			Password: "AlainPWD",
			Function: "volunteer",
		},
		"Sarah Croche": {
			Name:     "Sarah Croche",
			Password: "SarahPWD",
			Function: "organizer",
		},
		"Ondine Akeleur": {
			Name:     "Ondine Akeleur",
			Password: "OndinePWD",
			Function: "organizer",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestLoadEvents(t *testing.T) {

	got := LoadDatabaseFromJson(tests.GetTestData(t)).Events

	//want := map[uint]Event{
	//	0: {
	//		ID:        0,
	//		Name:      "Festival de la musique",
	//		Organizer: "Sarah Croche",
	//		Jobs: []Job{
	//			{
	//				ID:       0,
	//				Name:     "Buvette",
	//				Required: 2,
	//				Volunteers: []string{
	//					"Alex Terrieur",
	//				},
	//			},
	//			{
	//				ID:         1,
	//				Name:       "Sécurité",
	//				Required:   3,
	//				Volunteers: []string{},
	//			},
	//		},
	//	},
	//	1: {
	//		ID:        1,
	//		Name:      "Fête de la science",
	//		Organizer: "Ondine Akeleur",
	//		Jobs: []Job{
	//			{
	//				ID:         0,
	//				Name:       "Buvette",
	//				Required:   2,
	//				Volunteers: []string{},
	//			},
	//			{
	//				ID:         1,
	//				Name:       "Sécurité",
	//				Required:   3,
	//				Volunteers: []string{},
	//			},
	//		},
	//	},
	//}

	want := 0
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestCreateEvent(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	testDb, err := db.CreateEvent("Test", "Sarah Croche")
	if err != nil {
		t.Error(err)
	}
	got := testDb.Events
	//	want := map[uint]Event{
	//		0: {
	//			ID:        0,
	//			Name:      "Festival de la musique",
	//			Organizer: "Sarah Croche",
	//			Jobs: []Job{
	//				{
	//					ID:       0,
	//					Name:     "Buvette",
	//					Required: 2,
	//					Volunteers: []string{
	//						"Alex Terrieur",
	//					},
	//				},
	//				{
	//					ID:         1,
	//					Name:       "Sécurité",
	//					Required:   3,
	//					Volunteers: []string{},
	//				},
	//			},
	//		},
	//		1: {
	//			ID:        1,
	//			Name:      "Fête de la science",
	//			Organizer: "Ondine Akeleur",
	//			Jobs: []Job{
	//				{
	//					ID:         0,
	//					Name:       "Buvette",
	//					Required:   2,
	//					Volunteers: []string{},
	//				},
	//				{
	//					ID:         1,
	//					Name:       "Sécurité",
	//					Required:   3,
	//					Volunteers: []string{},
	//				},
	//			},
	//		},
	//		2: {
	//			ID:        2,
	//			Name:      "Test",
	//			Organizer: "Sarah Croche",
	//			Jobs:      []Job{},
	//		},
	//	}
	want := 0

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestCreateEventErrorIfIdExists(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.CreateEvent("Festival de la musique", "Sarah Croche")
	if !errors.Is(err, ErrorEventExists) {
		t.Error("Error not raised")
	}
}

func TestCreateEventErrorIfNameIsEmpty(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.CreateEvent("", "Sarah Croche")
	if !errors.Is(err, ErrorEventNameEmpty) {
		t.Error("Error not raised")
	}
}

func TestCreateEventErrorIfOrganizerDoesntExist(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.CreateEvent("Test", "Bob")
	if !errors.Is(err, ErrorUserNotFound) {
		t.Error("Error not raised")
	}
}

func TestCreateEventErrorIfOrganizerIsEmpty(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.CreateEvent("Test", "")
	if !errors.Is(err, ErrorOrganizerEmpty) {
		t.Error("Error not raised")
	}
}

func TestGetEvent(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	testDb, err := db.GetEventByName("Festival de la musique")
	if err != nil {
		t.Error(err)
	}
	got := testDb
	//	want := &Event{
	//		ID:        0,
	//		Name:      "Festival de la musique",
	//		Organizer: "Sarah Croche",
	//		Jobs: []Job{
	//			{
	//				ID:       0,
	//				Name:     "Buvette",
	//				Required: 2,
	//				Volunteers: []string{
	//					"Alex Terrieur",
	//				},
	//			},
	//			{
	//				ID:         1,
	//				Name:       "Sécurité",
	//				Required:   3,
	//				Volunteers: []string{},
	//			},
	//		},
	//	}
	want := 0
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetEventErrorIfIdDoesntExist(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.GetEventByName("Bob's party")
	if err == nil {
		t.Error("Error not raised")
	}
}
