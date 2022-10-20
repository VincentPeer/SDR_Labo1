package models

import (
	"SDR_Labo1/src/communication/server/tests"
	"errors"
	"reflect"
	"testing"
)

func TestGetVolunteer(t *testing.T) {
	job := Job{
		ID:         1,
		Name:       "Test",
		Required:   1,
		Volunteers: []string{"Alex"},
	}
	got, err := job.GetVolunteer("Alex")
	if err != nil {
		t.Error(err)
	}
	want := "Alex"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetVolunteerErrorIfVolunteerDoesntExist(t *testing.T) {
	job := Job{
		ID:         1,
		Name:       "Test",
		Required:   1,
		Volunteers: []string{"Alex"},
	}
	_, err := job.GetVolunteer("Test")
	if !errors.Is(err, ErrorVolunteerNotFound) {
		t.Error("Error not raised")
	}
}

func TestVolunteerErrorIfVolunteerEmpty(t *testing.T) {
	job := Job{
		ID:         1,
		Name:       "Test",
		Required:   1,
		Volunteers: []string{"Alex"},
	}
	_, err := job.GetVolunteer("")
	if !errors.Is(err, ErrorVolunteerEmpty) {
		t.Error("Error not raised")
	}
}

func TestAddVolunteer(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.Events[0].Jobs[0].AddVolunteer("Alain Terrieur")
	if err != nil {
		t.Error(err)
	}
	got := db.Events[0]
	//	want := Event{
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
	//					"Alain Terrieur",
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

func TestAddVolunteerErrorIfVolunteerAlreadyExists(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t)).Events
	_, err := db[0].Jobs[0].AddVolunteer("Alex Terrieur")
	if !errors.Is(err, ErrorVolunteerExists) {
		t.Error("Error not raised")
	}
}

func TestAddVolunteerErrorIfVolunteerEmpty(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t)).Events
	_, err := db[0].Jobs[0].AddVolunteer("")
	if !errors.Is(err, ErrorVolunteerEmpty) {
		t.Error("Error not raised")
	}
}

func TestAddVolunteerErrorIfVolunteerAboveMaximum(t *testing.T) {
	job := Job{
		ID:         1,
		Name:       "Test",
		Required:   1,
		Volunteers: []string{"Alex"},
	}
	_, err := job.AddVolunteer("Bob")
	if !errors.Is(err, ErrorVolunteerAboveMaximum) {
		t.Error("Error not raised")
	}
}

func TestJobToString(t *testing.T) {
	job := Job{
		ID:         1,
		Name:       "Test",
		Required:   1,
		Volunteers: []string{"Alex", "Bob"},
	}
	got := job.ToString()
	want := "1 | Test | 1 | Alex - Bob"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
