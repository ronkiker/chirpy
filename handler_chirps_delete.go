package main

import (
	"net/http"
	"strconv"

	"github.com/ronkiker/chirpy/blob/master/authenticate"
)

func (cfg *apiConfig) HandleChirpsDelete(w http.ResponseWriter, r *http.Request) {
	// get token for authentication
	token, err := authenticate.GetBearer(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// validate token
	userId, err := authenticate.ValidateJWT(token, cfg.JWT)
	if err != nil {
		// if there was an error in the validation
		RespondWithError(w, http.StatusUnauthorized, "Invalid token")
	}

	if userId == "" {
		// no user found with token, send 403 status
		w.WriteHeader(403)
	}

	// capture chirp ID from path
	chirpIDString := r.PathValue("chirpID")

	// convert to int
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	// convert user id to int
	usr, err := strconv.ParseInt(userId, 0, 0)

	result, err := cfg.DB.DeleteChirp(chirpID, int(usr))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
	}
	if result == "wrong user" {
		w.WriteHeader(403)
	}
	if result == "success" {
		w.WriteHeader(204)
	}
}
