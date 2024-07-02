package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ronkiker/chirpy/blob/master/authenticate"
)

func (cfg *apiConfig) HandlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}
	token, err := authenticate.GetBearer(r.Header)
	if err != nil {
		RespondWithError(w, 401, "couldn't find JWT")
		return
	}

	subject, err := authenticate.ValidateJWT(token, cfg.JWT)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		RespondWithError(w, 401, "couldn't decode request")
		return
	}

	hash, err := authenticate.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, 401, "couldn't hash password")
		return
	}

	userId, err := strconv.Atoi(subject)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	user, err := cfg.DB.UpdateUser(userId, params.Email, hash)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	//	ChirpyRed: user.ChirpyRed,
	RespondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}
