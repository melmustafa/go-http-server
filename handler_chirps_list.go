package main

import (
	"log"
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsList(w http.ResponseWriter, r *http.Request) {
	authorId, err := strconv.Atoi(r.URL.Query().Get("author_id"))
	if err != nil {
		authorId = 0
	}

	sortOrder := r.URL.Query().Get("sort")

	dbChirps, err := cfg.DB.ListChirps()
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		if dbChirp.AuthorId == authorId || authorId == 0 {
			chirps = append(chirps, Chirp{
				ID:       dbChirp.ID,
				Body:     dbChirp.Body,
				AuthorId: dbChirp.AuthorId,
			})
		}
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortOrder == "desc" {
			return chirps[i].ID > chirps[j].ID
		}
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}
