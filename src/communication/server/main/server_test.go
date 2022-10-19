package main

import (
	"SDR_Labo1/src/communication/server/models"
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
	want := []models.User{
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

	got := loadConfig(getTestData(t)).Events
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
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
