package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadUsers(t *testing.T) {

	path, err := filepath.Abs("./config_test.json")

	if err != nil {
		t.Error(err)
	}

	got := loadUsers(path)
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
