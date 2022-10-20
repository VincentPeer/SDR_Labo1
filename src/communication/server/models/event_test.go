package models

import (
	"SDR_Labo1/src/communication/server/tests"
	"errors"
	"reflect"
	"testing"
)

func TestCreateJob(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	got, err := db.GetEvent(0)
	if err != nil {
		t.Error(err)
	}
	_, err = got.CreateJob("Test", 3)
	if err != nil {
		t.Error(err)
	}
	want := 0
	//want := &Event{
	//	ID:        0,
	//	Name:      "Festival de la musique",
	//	Organizer: "Sarah Croche",
	//	Jobs: make(map[uint]*Job {
	//		0 : {
	//				ID:       0,
	//				Name:     "Buvette",
	//				Required: 2,
	//				Volunteers: []string{
	//					"Alex Terrieur",
	//				},
	//			},
	//			1: {
	//				ID:         1,
	//				Name:       "Sécurité",
	//				Required:   3,
	//				Volunteers: []string{},
	//			},
	//			2: {
	//				ID:         2,
	//				Name:       "Test",
	//				Required:   3,
	//				Volunteers: []string{},
	//			}})
	//		}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestCreateJobErrorIfIdExists(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	event, err := db.GetEvent(0)
	if err != nil {
		t.Error(err)
	}
	_, err = event.CreateJob("Buvette", 3)
	if !errors.Is(err, ErrorJobExists) {
		t.Error("Error not raised")
	}
}

func TestCreateJobErrorIfNameIsEmpty(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	event, err := db.GetEvent(0)
	if err != nil {
		t.Error(err)
	}
	_, err = event.CreateJob("", 3)
	if !errors.Is(err, ErrorJobNameEmpty) {
		t.Error("Error not raised")
	}
}

func TestGetJob(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	event, err := db.GetEvent(0)
	if err != nil {
		t.Error(err)
	}
	got, err := event.GetJob(0)
	if err != nil {
		t.Error(err)
	}
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
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	event, err := db.GetEvent(0)
	if err != nil {
		t.Error(err)
	}
	_, err = event.GetJobByName("Test")
	if !errors.Is(err, ErrorJobNotFound) {
		t.Error("Error not raised")
	}
}

func TestGetJobErrorIfNameIsEmpty(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	event, err := db.GetEvent(0)
	if err != nil {
		t.Error(err)
	}
	_, err = event.GetJobByName("")
	if !errors.Is(err, ErrorJobNameEmpty) {
		t.Error("Error not raised")
	}
}

func TestEventToString(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	event, err := db.GetEvent(0)
	if err != nil {
		t.Error(err)
	}
	got := event.ToString()
	want := "0 | Festival de la musique | Sarah Croche"

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetJobsAsStringArray(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	event, err := db.GetEvent(0)
	if err != nil {
		t.Error(err)
	}
	got := event.GetJobsAsStringArray()
	want := []string{
		"0 | Buvette | 2 | Alex Terrieur",
		"1 | Sécurité | 3 | ",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
