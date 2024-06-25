package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ronkiker/chirpy/blob/master/authenticate"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (cfg *apiConfig) HandlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	hashPassword, err := authenticate.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "couldn't hash password")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, hashPassword)
	if err != nil {
		if errors.Is(err, errors.New("resource does not exist")) {
			RespondWithError(w, http.StatusConflict, "User already exists")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	RespondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}
