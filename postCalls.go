package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) handlePostChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleaned)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, 201, Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
	})
}

func validateChirp(body string) (string, error) {
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	return profanityCheck(params.Body), nil
}

func profanityCheck(text string) string {
	wordArray := strings.Split(text, " ")
	s := NewSet()
	s.Add("kerfuffle")
	s.Add("sharbert")
	s.Add("fornax")
	for word := range wordArray {
		if s.Contains(strings.ToLower(wordArray[word])) {
			wordArray[word] = "****"
		}
	}
	return strings.Join(wordArray, " ")
}
