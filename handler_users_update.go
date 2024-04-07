package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	id, issuer, err := cfg.auth.ValidateToken(authToken)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	if issuer != cfg.auth.AccessIssuer {
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.UpdateUser(id, params.Email, params.Password)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:    user.ID,
		Email: user.Email,
	})
}
