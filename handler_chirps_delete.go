package main

import (
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	authorId, issuer, err := cfg.auth.ValidateToken(authToken)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	if issuer != cfg.auth.AccessIssuer {
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	id, err := strconv.Atoi(r.PathValue("chirpId"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse the id path parameter")
		return
	}

	deleted, found, err := cfg.DB.DeleteChirp(id, authorId)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't access the db")
		return
	}

	if !found {
		respondWithError(w, http.StatusNotFound, "")
		return
	}

	if found && !deleted {
		respondWithError(w, http.StatusForbidden, "")
		return
	}

	respondWithJSON(w, http.StatusOK, "")
}
