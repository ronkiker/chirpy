package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ronkiker/chirpy/blob/master/authenticate"
	"github.com/ronkiker/chirpy/blob/master/database"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
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

	token, err := authenticate.GetBearer(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	user, err := cfg.DB.GetUserForRefreshToken(token)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	var author int
	if (database.User{}) == user {
		userId, err := authenticate.ValidateJWT(token, cfg.JWT)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		}
		if len(userId) == 0 {
			w.WriteHeader(401)
			return
		}
		usr, err := strconv.ParseInt(userId, 0, 0)
		author = int(usr)
	} else {
		author = int(user.ID)
	}

	chirp, err := cfg.DB.CreateChirp(cleaned, int(user.ID))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}
	fmt.Printf("AUTHOR ID: %v \n", author)
	RespondWithJSON(w, 201, Chirp{
		ID:       chirp.ID,
		Body:     chirp.Body,
		AuthorId: author,
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
