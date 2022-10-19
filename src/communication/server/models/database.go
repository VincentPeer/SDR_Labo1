package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Database struct {
	Events Events `json:"events"`
	Users  Users  `json:"users"`
}

func LoadDatabaseFromJson(jsonPath string) Database {
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println("Error reading config file:", err.Error())
		os.Exit(1)
	}
	//fmt.Println("Successfully opened " + jsonFile.Name())
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var conf Database
	json.Unmarshal(byteValue, &conf)

	return conf
}
