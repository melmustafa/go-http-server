package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	users, err := cfg.DB.ListUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't connect to the db")
		return
	}

	for _, user := range users {
		if user.Email == params.Email && user.Password == params.Password {
			respondWithJSON(w, http.StatusOK, User{
				ID:    user.ID,
				Email: user.Email,
			})
			return
		}
	}

	respondWithError(w, http.StatusUnauthorized, "")
}
