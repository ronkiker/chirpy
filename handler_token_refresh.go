package main

import (
	"net/http"
	"time"

	"github.com/ronkiker/chirpy/blob/master/authenticate"
	"github.com/ronkiker/chirpy/blob/master/database"
)

func (cfg *apiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		RefreshToken string `json:"token"`
	}
	refreshToken, err := authenticate.GetBearer(r.Header)
	if err != nil {
		RespondWithError(w, 401, "Couldn't find token")
		return
	}

	user, err := cfg.DB.GetUserForRefreshToken(refreshToken)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	if (database.User{}) == user {
		w.WriteHeader(401)
		return
	}

	accessToken, err := authenticate.CreateJWT(
		user.ID,
		cfg.JWT,
		time.Hour,
	)
	if err != nil {
		RespondWithError(w, 401, "Unable to validate token")
		return
	}

	RespondWithJSON(w, http.StatusOK, response{
		RefreshToken: accessToken,
	})
}

func (cfg *apiConfig) HandlerRefreshRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := authenticate.GetBearer(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Couldn't find token")
	}
	err = cfg.DB.RevokeRefreshToken(refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Could not revoke session")
	}

	w.WriteHeader(http.StatusNoContent)
}
