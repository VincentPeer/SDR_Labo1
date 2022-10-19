package main

import (
	"SDR_Labo1/src/communication/server/models"
	"reflect"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db := loadConfig(getTestData(t)).Users
	testDb, err := models.CreateUser(db, "Test", "TestPWD", "volunteer")
	if err != nil {
		t.Error(err)
	}
	got := testDb
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
		{
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
	_, err := models.CreateUser(testDb, "Test", "TestPWD", "volunteer")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestGetUser(t *testing.T) {
	testDb := loadConfig(getTestData(t)).Users
	got, _ := models.GetUser(testDb, "Alex Terrieur")
	want := models.User{
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
	_, err := models.GetUser(testDb, "Test")
	if err == nil {
		t.Errorf("got %v want %v", err, "User not found")
	}
}

func TestLogin(t *testing.T) {
	users = loadConfig(getTestData(t)).Users

	got, err := models.Login(users, "Alex Terrieur", "AlexPWD")
	want := true

	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	got, err = models.Login(users, "Alex Terrieur", "AlexPWD2")
	want = false

	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestLoginError(t *testing.T) {
	users = loadConfig(getTestData(t)).Users

	_, err := models.Login(users, "Test", "AlexPWD")
	if err == nil {
		t.Error("Expected error")
	}
}
