package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

const (
	TEST_CONFIG_FILE_PATH = "./config_test.json"
)

func getTestData(t *testing.T) string {
	path, err := filepath.Abs(TEST_CONFIG_FILE_PATH)

	if err != nil {
		t.Error(err)
	}
	return path
}

func TestLoadUsers(t *testing.T) {

	got := loadConfig(getTestData(t)).Users
	want := []User{
		{
			Id:       "1",
			Name:     "Alex Terrieur",
			Password: "AlexPWD",
			Function: "volunteer",
		},
		{
			Id:       "2",
			Name:     "Alain Terrieur",
			Password: "AlainPWD",
			Function: "volunteer",
		},
		{
			Id:       "3",
			Name:     "Sarah Croche",
			Password: "SarahPWD",
			Function: "organizer",
		},
		{
			Id:       "4",
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

	got := loadConfig(getTestData(t)).Events
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
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestLogin(t *testing.T) {
	users = loadConfig(getTestData(t)).Users

	got := login("1", "AlexPWD")
	want := true

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	got = login("1", "AlexPWD2")
	want = false

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
