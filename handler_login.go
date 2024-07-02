package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ronkiker/chirpy/blob/master/authenticate"
)

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	//	ID        int    `json:"id"`
	//Email     string `json:"email"`
	//	ChirpyRed bool   `json:"is_chirpy_red"`
	//RefreshToken string `json:"refresh_token"`
	type response struct {
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
		RespondWithError(w, http.StatusInternalServerError, "Unable to find user")
		return
	}
	//chirpyRed := user.ChirpyRed
	err = authenticate.CheckPassword(params.Password, user.HashedPassword)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unable to find user")
		return
	}

	accessToken, err := authenticate.CreateJWT(
		user.ID,
		cfg.JWT,
		time.Hour,
	)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}
	refreshToken, err := authenticate.CreateRefreshToken()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token")
		return
	}

	err = cfg.DB.SaveRefreshToken(user.ID, refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token")
		return
	}
	//
	//ID:        user.ID,
	//Email:     user.Email,
	//ChirpyRed: chirpyRed,
	//RefreshToken: refreshToken,
	RespondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})

}
