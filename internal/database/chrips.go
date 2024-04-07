package database

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:       id,
		Body:     body,
		AuthorId: authorId,
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

func (db *DB) GetChirp(id int) (Chirp, bool, error) {
	db.mux.RLock()
	dbData, err := db.loadDB()
	db.mux.RUnlock()
	if err != nil {
		return Chirp{}, false, err
	}
	chirp, ok := dbData.Chirps[id]
	return chirp, ok, nil
}

func (db *DB) DeleteChirp(id, authorId int) (bool, bool, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()
	if err != nil {
		return false, false, err
	}

	chirp, ok := dbData.Chirps[id]
	if !ok {
		return false, false, nil
	}

	if chirp.AuthorId != authorId {
		return false, true, nil
	}

	delete(dbData.Chirps, id)
	db.writeDB(dbData)

	return true, true, nil
}
