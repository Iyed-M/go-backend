package database

import "log"

func (db *DB) AddRevokeToken(id string, token string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	dbStruct.RevokedTokens[id] = append(dbStruct.RevokedTokens[id], token)
	if err := db.writeDB(*dbStruct); err != nil {
		return err
	}

	return nil
}

func (db *DB) IsTokenRevoked(id string, token string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		log.Fatal("error loading db")
	}
	for _, revokedToken := range dbStruct.RevokedTokens[id] {
		if token == revokedToken {
			return true
		}
	}

	return false
}
