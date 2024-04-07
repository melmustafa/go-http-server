package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	users, err := cfg.DB.ListUsers()
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't connect to the db")
		return
	}

	for _, user := range users {
		if user.Email == params.Email && user.Password == params.Password {
			jwtToken, err := cfg.auth.CreateToken(user.ID, cfg.auth.AccessDuration, cfg.auth.AccessIssuer)
			if err != nil {
				log.Println(err)
				respondWithError(w, http.StatusInternalServerError, "couldn't generate jwt token")
				return
			}
			refreshToken, err := cfg.auth.CreateToken(user.ID, cfg.auth.RefreshDuration, cfg.auth.RefreshIssuer)
			if err != nil {
				log.Println(err)
				respondWithError(w, http.StatusInternalServerError, "couldn't generate refresh token")
				return
			}
			respondWithJSON(w, http.StatusOK, struct {
				Token        string `json:"token"`
				RefreshToken string `json:"refresh_token"`
				User
			}{
				Token:        jwtToken,
				RefreshToken: refreshToken,
				User: User{
					Email:       user.Email,
					ID:          user.ID,
					IsChirpyRed: user.IsChirpyRed,
				},
			})
			return
		}
	}

	respondWithError(w, http.StatusUnauthorized, "")
}
