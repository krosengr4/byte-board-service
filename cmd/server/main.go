package main

import (
	"byte-board/internal/appconfig"
	"byte-board/internal/handler"
	"byte-board/internal/middleware"
	"net/http"
	"os"
	"time"

	database "byte-board/internal/repository"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup Zerologger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05",
	}).
		With().
		Timestamp().
		Logger()

	log.Info().Msg("Starting Byte Board Backend Service!")

	// Load configurations
	cfg, err := appconfig.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize database
	db, err := database.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer db.Close()

	// Initialize handlers
	handler := handler.New(db, cfg)

	// Set up router with middlewear
	router := setupRouter(handler)

	// Initialize CORS middleware with configuration
	corsConfig := middleware.CORSConfig{
		AllowedOrigins: cfg.GetAllowedOrigins(),
	}

	// Apply middleware chain: Recover -> Logging -> CORS -> Router
	httpHandler := middleware.Recovery(
		middleware.Logging(
			middleware.CORS(corsConfig)(router),
		),
	)

	// Start server
	log.Info().Str("port", cfg.Port).Msg("Byte Board Service starting")

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      httpHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal().Err(server.ListenAndServe()).Msg("Server failed to start")
}

// Setup router configures all of the API routes
func setupRouter(h *handler.Handler) *mux.Router {
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Comments
	api.HandleFunc("/comments", h.GetAllComments).Methods("GET")
	api.HandleFunc("/post/{postId}/comments", h.GetCommentsOnPost).Methods("GET")

	return router
}
