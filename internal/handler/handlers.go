package handler

import (
	"byte-board/internal/appconfig"
	"byte-board/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
		return
	}

	log.Info().Int("count", len(comments)).Msg("Successfully retrieved comments!")
	writeJSONResponse(w, http.StatusOK, comments)
}

// Handler to get a comment by comment ID
func (h *Handler) GetCommentById(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET /comments/{CommentID} - Getting comment by its ID")

	vars := mux.Vars(r)
	idStr := vars["commentId"]

	log.Info().Str("comment_id", idStr).Msg("GET /comments/{CommentID} - Getting comment by ID")

	// Convert id string into an int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().Str("id", idStr).Msg("Invalid commend ID format")
	}

	// Get comment by id from the database
	comment, err := h.db.GetCommentById(id)
	if err != nil {
		if err.Error() == "comment not found" {
			log.Warn().Int("ID", id).Msg("Comment with that ID not found")
			writeErrorResponse(w, http.StatusInternalServerError, "Comment not found")
			return
		}
		log.Error().Err(err).Msg("Failed to get comment by ID")
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get that comment")
		return
	}

	log.Info().Int("ID", id).Msg("Successfully retrieved the comment")
	writeJSONResponse(w, http.StatusOK, comment)
}

// Handler to get all of the comments on a post
func (h *Handler) GetCommentsOnPost(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET /post/{postId}/comments - Getting comments on post")

	vars := mux.Vars(r)
	idStr := vars["postId"]

	// Convert the ID string into an int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().Str("id", idStr).Msg("Invalid post ID format")
		writeErrorResponse(w, http.StatusBadRequest, "Invalid Post ID")
		return
	}

	comments, err := h.db.GetCommentsByPost(id)
	if err != nil {
		log.Error().Err(err).Msg("GET /post/{postId}/comments - Getting all comments on a post")
		writeErrorResponse(w, http.StatusInternalServerError, "failed to get comments on post")
		return
	}

	log.Info().Int("count", len(comments)).Msg("Successfully retrieved comments on post")
	writeJSONResponse(w, http.StatusOK, comments)

}

// #endregion
