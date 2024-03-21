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
	Chrips map[string]Chirp `json:"chirps"`
}

func NewDB(path string) *DB {
	db := &DB{
		id:   0,
		path: path,
	}

	db.ensureDB()

	return db
}

func (db *DB) getLastID() int {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.id++
	return db.id
}

// CreateChirp creates a new chirp from body and return it
func (db *DB) CreateChirp(body string) Chirp {
	chirp := Chirp{
		ID:   db.getLastID(),
		Body: body,
	}

	// load db in memory
	db.mu.Lock()
	defer db.mu.Unlock()

	dbstruct, err := db.loadDB()
	if err != nil && err.Error() != "empty file" {
		return Chirp{}
	}

	dbstruct.Chrips[strconv.Itoa(chirp.ID)] = chirp

	err = db.writeDB(*dbstruct)
	if err != nil {
		return Chirp{}
	}

	return chirp
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil && err.Error() != "empty file" {
		return nil, err
	}
	chirps := []Chirp{}

	for _, chirp := range (*dbStruct).Chrips {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}
