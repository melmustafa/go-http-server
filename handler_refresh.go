package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	id, issuer, err := cfg.auth.ValidateToken(authToken)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}
	if issuer != cfg.auth.RefreshIssuer {
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	ok, err := cfg.DB.IsValidToken(authToken)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't validate your token")
		return
	}

	if !ok {
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	token, err := cfg.auth.CreateToken(id, cfg.auth.AccessDuration, cfg.auth.AccessIssuer)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate token")
		return
	}
	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}
