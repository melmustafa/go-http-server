package database

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) ListChirps() ([]Chirp, error) {
	db.mux.RLock()
	dbData, err := db.loadDB()
	db.mux.RUnlock()
	if err != nil {
		return nil, err
	}
	chirps := make([]Chirp, 0)
	for _, chirp := range dbData.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) GetChirps(id int) (Chirp, bool, error) {
	db.mux.RLock()
	dbData, err := db.loadDB()
	db.mux.RUnlock()
	if err != nil {
		return Chirp{}, false, err
	}
	chirp, ok := dbData.Chirps[id]
	return chirp, ok, nil
}
