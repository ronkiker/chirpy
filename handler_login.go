package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ronkiker/chirpy/blob/master/authenticate"
)

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = authenticate.CheckPassword(params.Password, user.HashedPassword)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}
	defaultExp := 60 * 60 * 24
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = defaultExp
	} else if params.ExpiresInSeconds > defaultExp {
		params.ExpiresInSeconds = defaultExp
	}

	token, err := authenticate.CreateJWT(user.ID, cfg.JWT, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create token")
		return
	}
	RespondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
		Token: token,
	})

}
