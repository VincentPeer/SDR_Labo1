package main

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Function string `json:"function"`
}
