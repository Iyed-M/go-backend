package database

import "strconv"

type Chirp struct {
	Body string `json:"body"`
	ID   int    `json:"id"`
}

// CreateChirp creates a new chirp from body and return it
func (db *DB) CreateChirp(body string) (Chirp, error) {
	chirp := Chirp{
		ID:   db.getLastChirpID(),
		Body: body,
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
