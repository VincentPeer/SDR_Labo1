package main

import (
	"reflect"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db := loadConfig(getTestData(t)).Users
	testDb, err := createUser(db, "5", "Test", "TestPWD", "volunteer")
	if err != nil {
		t.Error(err)
	}
	got := testDb
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
		{
			Id:       "5",
			Name:     "Test",
			Password: "TestPWD",
			Function: "volunteer",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestCreateUserError(t *testing.T) {
	testDb := loadConfig(getTestData(t)).Users
	_, err := createUser(testDb, "1", "Test", "TestPWD", "volunteer")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestGetUser(t *testing.T) {
	testDb := loadConfig(getTestData(t)).Users
	got, _ := getUser(testDb, "1")
	want := User{
		Id:       "1",
		Name:     "Alex Terrieur",
		Password: "AlexPWD",
		Function: "volunteer",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetUserError(t *testing.T) {
	testDb := loadConfig(getTestData(t)).Users
	_, err := getUser(testDb, "5")
	if err == nil {
		t.Errorf("got %v want %v", err, "User not found")
	}
}
