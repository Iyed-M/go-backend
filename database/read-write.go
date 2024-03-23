package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var ErrEmptyFile = errors.New("empty file")

func newDBStructure() *DBStructure {
	return &DBStructure{
		Chrips:        make(map[string]Chirp),
		Users:         make(map[string]User),
		RevokedTokens: make(map[string][]string),
	}
}

// newId generares new ID by incrementing it each time it s called
func (db *DB) newId() int {
	db.chirpID++
	return db.chirpID
}

// loadDB loads the databasefile into memory returns a pointer to it
// throw erro empty file if the file is empty
func (db *DB) loadDB() (*DBStructure, error) {
	if err := db.ensureDB(); err != nil {
		return &DBStructure{}, err
	}
	dataJSON, err := os.ReadFile(db.path)
	if err != nil {
		return &DBStructure{}, err
	}
	if len(dataJSON) == 0 {
		return newDBStructure(), ErrEmptyFile
	}

	buffer := newDBStructure()

	err = json.Unmarshal(dataJSON, &buffer)
	if err != nil {
		return &DBStructure{}, err
	}

	return buffer, nil
}

// writeDB write db file to disk
func (db *DB) writeDB(bufferDB DBStructure) error {
	dataJSON, err := json.Marshal(bufferDB)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dataJSON, 0o644)
	if err != nil {
		return err
	}

	return nil
}

// ensureDB checks if the database file exists and creates it if it doesn't
func (db *DB) ensureDB() error {
	// check if db file db.json exists
	_, err := os.ReadFile(db.path)

	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		// create db file
		f, err := os.Create("database.json")
		f.Close()
		if err != nil {
			return fmt.Errorf("error createating file : %v", err)
		}
	}

	return nil
}
