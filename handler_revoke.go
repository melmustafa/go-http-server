package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	_, issuer, err := cfg.auth.ValidateToken(authToken)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't validate your token")
		return
	}

	if issuer != cfg.auth.RefreshIssuer {
		respondWithError(w, http.StatusBadRequest, "can only revoke refresh tokens")
	}

	ok, err := cfg.DB.RevokeToken(authToken)

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke your token")
		return
	}

	if !ok {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke your token")
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
}
