package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey := strings.TrimPrefix(r.Header.Get("Authorization"), "ApiKey ")
	if apiKey != cfg.apiKey {
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't decode request body")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, "")
		return
	}

	err = cfg.DB.UpgradeUser(params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't access the db")
		return
	}

	respondWithJSON(w, http.StatusOK, "")
}
