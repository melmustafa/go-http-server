package database

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red "`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:          id,
		Email:       email,
		Password:    password,
		IsChirpyRed: false,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) ListUsers() ([]User, error) {
	db.mux.RLock()
	dbData, err := db.loadDB()
	db.mux.RUnlock()
	if err != nil {
		return nil, err
	}
	users := make([]User, 0)
	for _, user := range dbData.Users {
		users = append(users, user)
	}
	return users, nil
}

func (db *DB) UpdateUser(id int, email, password string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbData, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	originalUser := dbData.Users[id]

	user := User{
		ID:          id,
		Email:       email,
		Password:    password,
		IsChirpyRed: originalUser.IsChirpyRed,
	}

	dbData.Users[id] = user

	err = db.writeDB(dbData)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpgradeUser(id int) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbData, err := db.loadDB()
	if err != nil {
		return err
	}

	originalUser := dbData.Users[id]

	originalUser.IsChirpyRed = true

	dbData.Users[id] = originalUser

	err = db.writeDB(dbData)
	if err != nil {
		return err
	}

	return nil
}
