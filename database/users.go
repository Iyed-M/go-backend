package database

import (
	"strconv"
)

type User struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}

// CreateUser creates a new user from email , stores it and returns it
func (db *DB) CreateUser(email string) (User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	usr := User{
		ID:    db.getLastUserID(),
		Email: email,
	}

	dbStruct, err := db.loadDB()
	if err != nil && err.Error() != "empty file" {
		return User{}, err
	}

	dbStruct.Users[strconv.Itoa(usr.ID)] = usr

	db.writeDB(*dbStruct)

	return usr, nil
}
