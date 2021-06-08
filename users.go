package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (user *User) addToDataStore() error {
	users, err := getAllUsers()
	if err != nil {
		return err
	}

	lastUserID := users[len(users)-1].ID
	user.ID = lastUserID + 1
	users = append(users, *user)

	err = saveUsersToDataStore(users)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) deleteFromDataStore() error {
	users, err := getAllUsers()
	if err != nil {
		return err
	}

	for i, searchedUser := range users {
		if user.ID == searchedUser.ID {

			users = append(users[:i], users[i+1:]...)
			err = saveUsersToDataStore(users)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func (user *User) updateToDataStore() error {
	users, err := getAllUsers()
	if err != nil {
		return err
	}

	for i, searchedUser := range users {
		if user.ID == searchedUser.ID {
			users[i] = *user
			err = saveUsersToDataStore(users)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("error updating the user")
}

func (user *User) validate() string {
	if user.Name == "" {
		return "username cannot be empty"
	}

	return ""
}

func getAllUsers() ([]User, error) {
	var users []User

	jsonFile, err := os.Open("data/users.json")
	defer jsonFile.Close()
	if err != nil {
		return users, err
	}

	jsonBytes, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(jsonBytes, &users)

	return users, nil
}

func getUserByID(id int) (User, error) {
	var notFoundUser User
	notFoundUser.ID = -1

	users, err := getAllUsers()
	if err != nil {
		return notFoundUser, err
	}

	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return notFoundUser, nil
}

func saveUsersToDataStore(users []User) error {
	jsonData, err := json.MarshalIndent(users, "", "	")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("data/users.json", jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
