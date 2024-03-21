package database

import (
	"sync"
)

type DB struct {
	path    string
	chirpID int
	userID  int
	mu      sync.Mutex
}

type DBStructure struct {
	Chrips map[string]Chirp `json:"chirps"`
	Users  map[string]User  `json:"users"`
}

func NewDB(path string) *DB {
	db := &DB{
		chirpID: 0,
		userID:  0,
		path:    path,
	}

	db.ensureDB()

	return db
}

// getLastUserID increments the ID of the latest chirp created and returns it
func (db *DB) getLastUserID() int {
	db.userID++
	return db.userID
}

// getLastChirpID increments the ID of the latest chirp created and returns it
func (db *DB) getLastChirpID() int {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.chirpID++
	return db.chirpID
}
