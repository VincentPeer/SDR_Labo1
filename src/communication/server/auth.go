package main

func login(id string, password string) bool {
	for _, user := range users {
		if user.Id == id && user.Password == password {
			return true
		}
	}
	return false
}
