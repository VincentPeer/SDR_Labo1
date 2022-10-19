package models

import (
	"SDR_Labo1/src/communication/server/tests"
	"reflect"
	"testing"
)

func TestLoadUsers(t *testing.T) {

	got := LoadDatabaseFromJson(tests.GetTestData(t)).Users
	want := Users{
		{
			Name:     "Alex Terrieur",
			Password: "AlexPWD",
			Function: "volunteer",
		},
		{
			Name:     "Alain Terrieur",
			Password: "AlainPWD",
			Function: "volunteer",
		},
		{
			Name:     "Sarah Croche",
			Password: "SarahPWD",
			Function: "organizer",
		},
		{
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
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
