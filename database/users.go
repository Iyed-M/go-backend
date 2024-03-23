package database

import (
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrIncorrectPassword  = errors.New("incorrect password")
	ErrEmailNotFound      = errors.New("email not found")
)

type User struct {
	Email string `json:"email"`
	// Password is the hashed password of the user.
	Password string `json:"password"`
	ID       int    `json:"id"`
}

// CreateUser creates a new user from email, password , stores it and returns it .
//
// email must be untique.
//
// error EmailAlreadyExists is returned if the email already exists.
func (db *DB) CreateUser(email string, password string) (User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil && err.Error() != "empty file" {
		return User{}, err
	}

	if !isEmailUnique(email, dbStruct.Users) {
		return User{}, ErrEmailAlreadyExists
	}

	usr := User{
		ID:       db.getLastUserID(),
		Email:    email,
		Password: password,
	}
	dbStruct.Users[strconv.Itoa(usr.ID)] = usr
	db.writeDB(*dbStruct)

	return usr, nil
}

func (db *DB) GetUserIDByEmail(email string, password string) (id int, err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	dbStruct, err := db.loadDB()
	if err != nil && err.Error() != "empty file" {
		return 0, err
	}

	user := User{}
	for _, user = range dbStruct.Users {
		if user.Email == email {
			err = validateUserPassword(password, user.Password)
			if err != nil {
				return 0, err
			}
			return user.ID, nil
		}
	}
	return 0, ErrEmailNotFound
}

// validateUserPassword validates the password of a user.
//
// returns ErrIncorrectPassword if the password is incorrect.
// returns nil if the password is correct.
func validateUserPassword(password string, hashedPassword string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return ErrIncorrectPassword
	}
	return nil
}

func (db *DB) UpdateUser(newUser User) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	loadedDB, err := db.loadDB()
	if err != nil {
		return err
	}

	fmt.Println("new password:", loadedDB.Users[strconv.Itoa(newUser.ID)].Password)
	fmt.Println("new password:", newUser.Password)
	loadedDB.Users[strconv.Itoa(newUser.ID)] = newUser
	db.writeDB(*loadedDB)
	return nil
}
