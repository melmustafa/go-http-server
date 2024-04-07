package database

import (
	"strings"
	"time"
)

func (db *DB) RevokeToken(authToken string) (bool, error) {
	token := strings.TrimPrefix(authToken, "Bearer ")
	db.mux.Lock()
	defer db.mux.Unlock()
	dbData, err := db.loadDB()
	if err != nil {
		return false, err
	}
	_, ok := dbData.Tokens[token]
	if ok {
		return false, nil
	}
	dbData.Tokens[token] = time.Now()
	err = db.writeDB(dbData)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *DB) IsValidToken(authToken string) (bool, error) {
	token := strings.TrimPrefix(authToken, "Bearer ")
	db.mux.RLock()
	dbData, err := db.loadDB()
	defer db.mux.RUnlock()
	if err != nil {
		return false, err
	}
	_, ok := dbData.Tokens[token]
	return !ok, nil
}
