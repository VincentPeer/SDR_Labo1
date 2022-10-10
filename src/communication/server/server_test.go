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

	got := loadUsers(getTestData(t))
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
			Function: "organiser",
		},
		{
			Id:       "4",
			Name:     "Ondine Akeleur",
			Password: "OndinePWD",
			Function: "organiser",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestLogin(t *testing.T) {
	users = loadUsers(getTestData(t))

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
