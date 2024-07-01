package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/ronkiker/chirpy/blob/master/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	JWT            string
	POLKA          string
}

func main() {
	const root = "."
	const port = "8080"

	err := godotenv.Load("variables.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET not set")
	}

	db, err := database.NewDB("database/database.json")
	if err != nil {
		log.Fatal(err)
	}

	debug := flag.Bool("debug", false, "Enable debug")
	flag.Parse()
	if debug != nil && *debug {
		err := db.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		JWT:            jwtSecret,
		POLKA:          polkaKey,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(root)))))

	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/login", apiCfg.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandlerRefreshRevoke)

	mux.HandleFunc("PUT /api/users", apiCfg.HandlerUsersUpdate)
	mux.HandleFunc("POST /api/users", apiCfg.HandlerUserCreate)

	mux.HandleFunc("POST /api/chirps", apiCfg.HandleChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandlerChirpsRetrieve)

	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandlerChirpsGet)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.HandleChirpsDelete)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.HandlePolkaPost)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("listening for files from %v on port %v ", root, port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits)))
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Repsonsing with 5XX error: %s \n", msg)
	}
	if code == 404 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
	} else {
		type errorResponse struct {
			Error string `json:"error"`
		}
		RespondWithJSON(w, code, errorResponse{
			Error: msg,
		})
	}
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error Marshalling JSON: %s \n", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
