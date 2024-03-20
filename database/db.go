package database

import (
	"strconv"
	"sync"
)

type DB struct {
	path string
	id   int
	mu   sync.Mutex
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type Chirps []Chirp

type DBStructure struct {
	Chrips map[string]Chirp
}

func NewDB(path string) *DB {
	db := &DB{
		id:   1,
		path: path,
	}

	db.ensureDB()

	return db
}

// CreateChirp creates a new chirp and adds it to the database
func (db *DB) CreateChirp(chirp Chirp) error {
	// Stores all chirps in the memory
	db.mu.Lock()
	defer db.mu.Unlock()
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	// add new chirp in memory
	newID := db.newId()
	dbStruct.Chrips[strconv.Itoa(newID)] = chirp
	// write chrips in disk
	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := []Chirp{}
	for _, chirp := range dbStruct.Chrips {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}
