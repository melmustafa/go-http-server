package main

import (
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("chirpId"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse the id path parameter")
		return
	}

	chirp, found, err := cfg.DB.GetChirp(id)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't access the db")
		return
	}

	if !found {
		respondWithError(w, http.StatusNotFound, "")
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
