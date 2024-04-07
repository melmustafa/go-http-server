package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/melmustafa/go-http-server/internal/auth"
	"github.com/melmustafa/go-http-server/internal/database"
)

type apiConfig struct {
	fileserverHits int
	apiKey         string
	DB             *database.DB
	auth           auth.JWTConfig
}

func main() {
	const filepathRoot = "./pages"
	const port = "8080"

	godotenv.Load()

	db, err := database.NewDB("db.json")
	if err != nil {
		log.Fatal(err)
	}

	jwt := auth.JWTConfig{
		Secret:          os.Getenv("JWT_SECRET"),
		AccessIssuer:    "chirpy-access",
		RefreshIssuer:   "chirpy-refresh",
		AccessDuration:  time.Hour,
		RefreshDuration: time.Duration(60 * 24 * time.Hour),
	}

	apiCfg := apiConfig{
		apiKey:         os.Getenv("POLKA_API_KEY"),
		fileserverHits: 0,
		DB:             db,
		auth:           jwt,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsList)
	mux.HandleFunc("GET /api/chirps/{chirpId}", apiCfg.handlerChirpsGet)
	mux.HandleFunc("DELETE /api/chirps/{chirpId}", apiCfg.handlerChirpsDelete)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerUsersLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhook)

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
