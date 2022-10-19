package models

import "errors"

var (
	ErrorUserNotFound  = errors.New("user not found")
	ErrorUserExists    = errors.New("user with same id already exists")
	ErrorUserNameEmpty = errors.New("user name cannot be empty")
	ErrorPasswordEmpty = errors.New("password cannot be empty")
	ErrorFunctionEmpty = errors.New("function cannot be empty")
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Function string `json:"function"`
}

type Users []User

func (users *Users) ToMap() map[string]User {
	usersMap := make(map[string]User)
	for _, user := range *users {
		usersMap[user.Name] = user
	}
	return usersMap
}

// Creates a new user and adds it to the database
// Returns an error if a user with the same id already exists
// Otherwise returns the new state of the database
func (db *Database) CreateUser(name string, password string, function string) (*Database, error) {
	if name == "" {
		return nil, ErrorUserNameEmpty
	}
	if password == "" {
		return nil, ErrorPasswordEmpty
	}
	if function == "" {
		return nil, ErrorFunctionEmpty
	}
	if _, err := db.GetUser(name); err == nil {
		return nil, ErrorUserExists
	}
	db.Users[name] = User{name, password, function}
	return db, nil
}

// Get a user from the database
// Returns an error if the user does not exist
// Otherwise returns the user
func (db *Database) GetUser(name string) (User, error) {
	if name == "" {
		return User{}, ErrorUserNameEmpty
	}
	if user, ok := db.Users[name]; ok {
		return user, nil
	}
	return User{}, ErrorUserNotFound
}

// Confirm the password of a user
// Returns an error if the user does not exist
func (db Database) Login(name string, password string) (bool, error) {
	user, err := db.GetUser(name)
	if err != nil {
		return false, err
	}
	return user.Password == password, nil
}
