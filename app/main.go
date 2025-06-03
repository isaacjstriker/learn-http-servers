package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/isaacjstriker/learn-http-servers/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error fetching database: %s", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"

	cfg := apiConfig{
		dbQueries: dbQueries,
		platform:  os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()
	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", cfg.middlewareMetricsInc(fileHandler))

	// Endpoints
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", http.HandlerFunc(cfg.handlerMetrics))
	mux.HandleFunc("POST /admin/reset", http.HandlerFunc(cfg.handlerReset))
	mux.HandleFunc("POST /api/chirps", http.HandlerFunc(cfg.handlerCreateChirp))
	mux.HandleFunc("POST /api/users", http.HandlerFunc(cfg.handlerCreateUser))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// Server hits logic

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	hits := cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", hits)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden access")
		return
	}

	// Delete all users
	err := cfg.dbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Hits: 0"))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := map[string]string{"error": msg}
	data, _ := json.Marshal(resp)
	w.Write(data)
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	data, _ := json.Marshal(payload)
	w.Write(data)
}

func cleanProfanity(input string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(input, " ")
	for i, word := range words {
		lower := strings.ToLower(word)
		if _, found := badWords[lower]; found {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}
	type response struct {
        ID        string `json:"id"`
        Body      string `json:"body"`
        UserID    string `json:"user_id"`
        CreatedAt string `json:"created_at"`
        UpdatedAt string `json:"updated_at"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	// Clean profanity
	cleanedBody := cleanProfanity(req.Body)

	// Optionally: validate body length, user_id format, etc.
	userUUID, err := uuid.Parse(req.UserID)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid user_id format")
        return
    }

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
    Body:   cleanedBody,
    UserID: userUUID,
})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp")
		return
	}

	resp := response{
        ID:        chirp.ID.String(),
        Body:      chirp.Body,
        UserID:    req.UserID,
        CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
        UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
	}
	respondWithJSON(w, http.StatusCreated, resp)
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email string `json:"email"`
	}
	type response struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		log.Printf("Error decoding parameters: %s\n", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user")
		log.Printf("Error creating user: %s\n", err)
		return
	}

	resp := response{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusCreated, resp)
}
