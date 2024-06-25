package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) HandleChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleaned)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	RespondWithJSON(w, 201, Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
	})
}

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", errors.New("Chirp is too long")
	}
	return profanityCheck(body), nil
}

func profanityCheck(text string) string {
	wordsMap := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(text, wordsMap)
	return cleaned
}

func getCleanedBody(text string, wordsMap map[string]struct{}) string {
	words := strings.Split(text, " ")
	for x, word := range words {
		lowerCase := strings.ToLower(word)
		if _, ok := wordsMap[lowerCase]; ok {
			words[x] = "****"
		}
	}
	return strings.Join(words, " ")
}
