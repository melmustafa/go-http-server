package database

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
		ID:       id,
		Email:    email,
		Password: password,
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
