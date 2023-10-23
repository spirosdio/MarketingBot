package main

import "fmt"

type User struct {
	Name    string
	Surname string
	UserID  int
	Age     string
}

type UserDatabase struct {
	users map[int]*User
}

func NewUserDatabase() *UserDatabase {
	return &UserDatabase{
		users: make(map[int]*User),
	}
}

func (db *UserDatabase) CreateUser(name string, surname string, userID int, age string) *User {
	if _, exists := db.users[userID]; !exists {
		newUser := &User{
			Name:    name,
			Surname: surname,
			UserID:  userID,
			Age:     age,
		}
		db.users[userID] = newUser
		return newUser
	}
	return nil
}

func (db *UserDatabase) GetUser(userID int) *User {
	return db.users[userID]
}

func (db *UserDatabase) GetAllUsers() []*User {
	var allUsers []*User
	for _, user := range db.users {
		allUsers = append(allUsers, user)
	}
	return allUsers
}

func (db *UserDatabase) ListUsers() {
	for _, user := range db.users {
		// Print user details
		fmt.Println(user.Name, user.Surname, user.UserID, user.Age)
	}
}
