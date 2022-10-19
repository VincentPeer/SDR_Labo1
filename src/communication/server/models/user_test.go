package models

import (
	"SDR_Labo1/src/communication/server/tests"
	"reflect"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := db.CreateUser("Test", "TestPWD", "volunteer")
	if err != nil {
		t.Error(err)
	}
	got := db.Users
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

func TestCreateUserErrorIfUserAlreadyExist(t *testing.T) {
	testDb := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := testDb.CreateUser("Alex Terrieur", "TestPWD", "volunteer")
	if err != ErrorUserExists {
		t.Error("Expected error")
	}
}

func TestCreateUserErrorIfUserNameEmpty(t *testing.T) {
	testDb := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := testDb.CreateUser("", "TestPWD", "volunteer")
	if err != ErrorUserNameEmpty {
		t.Error("Expected error")
	}
}

func TestCreateUserErrorIfPasswordEmpty(t *testing.T) {
	testDb := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := testDb.CreateUser("Test", "", "volunteer")
	if err != ErrorPasswordEmpty {
		t.Error("Expected error")
	}
}

func TestCreateUserErrorIfFunctionEmpty(t *testing.T) {
	testDb := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := testDb.CreateUser("Test", "TestPWD", "")
	if err != ErrorFunctionEmpty {
		t.Error("Expected error")
	}
}

func TestGetUser(t *testing.T) {
	testDb := LoadDatabaseFromJson(tests.GetTestData(t))
	got, _ := testDb.GetUser("Alex Terrieur")
	want := User{
		Name:     "Alex Terrieur",
		Password: "AlexPWD",
		Function: "volunteer",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetUserError(t *testing.T) {
	testDb := LoadDatabaseFromJson(tests.GetTestData(t))
	_, err := testDb.GetUser("Test")
	if err == nil {
		t.Errorf("got %v want %v", err, "User not found")
	}
}

func TestLogin(t *testing.T) {
	users := LoadDatabaseFromJson(tests.GetTestData(t))
	got, err := users.Login("Alex Terrieur", "AlexPWD")
	want := true

	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	got, err = users.Login("Alex Terrieur", "AlexPWD2")
	want = false

	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestLoginError(t *testing.T) {
	users := LoadDatabaseFromJson(tests.GetTestData(t))

	_, err := users.Login("Test", "AlexPWD")
	if err == nil {
		t.Error("Expected error")
	}
}
