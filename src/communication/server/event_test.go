package main

import (
	"reflect"
	"testing"
)

func TestCreateEvent(t *testing.T) {
	db := loadConfig(getTestData(t)).Events
	testDb, err := createEvent(db, "3", "Test", "3")
	if err != nil {
		t.Error(err)
	}
	got := testDb
	want := []Event{
		{
			Id:        "1",
			Name:      "Festival de la musique",
			Organizer: "3",
			Jobs: []Job{
				{
					Id:       "1",
					Name:     "Buvette",
					Required: 2,
					Volunteers: []string{
						"1",
					},
				},
				{
					Id:         "2",
					Name:       "Sécurité",
					Required:   3,
					Volunteers: []string{},
				},
			},
		},
		{
			Id:        "2",
			Name:      "Fête de la science",
			Organizer: "4",
			Jobs: []Job{
				{
					Id:         "3",
					Name:       "Buvette",
					Required:   2,
					Volunteers: []string{},
				},
				{
					Id:         "4",
					Name:       "Sécurité",
					Required:   3,
					Volunteers: []string{},
				},
			},
		},
		{
			Id:        "3",
			Name:      "Test",
			Organizer: "3",
			Jobs:      []Job{},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
