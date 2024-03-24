package database

import (
	"errors"
	"strconv"
)

var ErrChirpNotFound = errors.New("chirp not found")

type Chirp struct {
	Body     string `json:"body"`
	ID       int    `json:"id"`
	AuthorId string `json:"author_id"`
}

// CreateChirp creates a new chirp from body and return it
func (db *DB) CreateChirp(authId, body string) (Chirp, error) {
	chirp := Chirp{
		ID:       db.getLastChirpID(),
		Body:     body,
		AuthorId: authId,
	}

	// load db in memory
	db.mu.Lock()
	defer db.mu.Unlock()

	dbstruct, err := db.loadDB()
	if err != nil && err.Error() != "empty file" {
		return Chirp{}, err
	}

	dbstruct.Chrips[strconv.Itoa(chirp.ID)] = chirp

	err = db.writeDB(*dbstruct)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// GetChirps returns all chirps
// error "empty file" is returned if no chrips are found in db
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

func (db *DB) DeleteChirp(id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	delete(dbStruct.Chrips, id)
	db.writeDB(*dbStruct)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetChirpByID(chirpID string) (Chirp, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	for id, chirp := range dbStruct.Chrips {
		if id == chirpID {
			return chirp, nil
		}
	}

	return Chirp{}, ErrChirpNotFound
}
