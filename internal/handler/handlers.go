package handler

import (
	"byte-board/internal/appconfig"
	"byte-board/internal/repository"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Handler struct {
	db     *repository.DB
	config *appconfig.Config
}

// Create a new instance of a handler
func New(db *repository.DB, cfg *appconfig.Config) *Handler {
	return &Handler{
		db:     db,
		config: cfg,
	}
}

// Represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Writes a JSON response
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error().Err(err).Msg("Error encoding JSON response")
	}
}

// Writes an error response
func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	log.Warn().Int("status", status).Str("message", message).Msg("Writing error response")
	writeJSONResponse(w, status, ErrorResponse{Error: message})
}

// #region Comments

// Handler to get all comments
func (h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET /comments - Getting all comments")

	comments, err := h.db.GetAllComments()
	if err != nil {
		log.Error().Err(err).Msg("Error getting comments")
		writeErrorResponse(w, http.StatusInternalServerError, "failed to get comments")
	}

	log.Info().Int("count", len(comments)).Msg("Successfully retrieved comments!")
	writeJSONResponse(w, http.StatusOK, comments)
}

// #endregion
