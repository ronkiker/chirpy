package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) HandlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid chirp")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	})
}

func (cfg *apiConfig) HandlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	sorted := r.URL.Query().Get("sort")
	author := r.URL.Query().Get("author_id")

	if len(sorted) != 0 && sorted != "asc" && sorted != "desc" {
		RespondWithError(w, http.StatusBadRequest, "Invalid sort")
	}

	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error retrieving chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		if len(author) > 0 {
			author, err := strconv.Atoi(author)
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, "bad id")
			}
			if dbChirp.AuthorId == author {
				chirps = append(chirps, Chirp{
					ID:   dbChirp.ID,
					Body: dbChirp.Body,
				})
			}
		} else {
			chirps = append(chirps, Chirp{
				ID:   dbChirp.ID,
				Body: dbChirp.Body,
			})
		}

	}

	sort.Slice(chirps, func(i, j int) bool {
		if len(sorted) == 0 || sorted == "asc" {
			return chirps[i].ID > chirps[j].ID
		} else {
			return chirps[i].ID < chirps[j].ID
		}

	})
	RespondWithJSON(w, http.StatusOK, chirps)
}
