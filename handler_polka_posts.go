package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/ronkiker/chirpy/blob/master/authenticate"
)

func (cfg *apiConfig) HandlePolkaPost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	polkaKey, err := authenticate.GetApiKey(r.Header)
	if err != nil {
		w.WriteHeader(401)
	}
	if len(polkaKey) == 0 {
		w.WriteHeader(401)
	}
	if polkaKey != cfg.POLKA {
		w.WriteHeader(401)
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	user, err := cfg.DB.UpdateUserToRed(params.Data.UserID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	if errors.Is(err, os.ErrNotExist) {
		w.WriteHeader(404)
	}
	if user == "success" {
		w.WriteHeader(204)
	}

}
