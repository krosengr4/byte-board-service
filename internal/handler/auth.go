package handler

import (
	"byte-board/internal/model"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// POST /api/register - Register handler
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("POST /api/register - Registering new user")

	// Parse body request
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		log.Warn().Msg("Missing required fields")
		writeErrorResponse(w, http.StatusBadRequest, "Username, password, first name, and last name are required")
		return
	}

	// Create user and profile with auth service
	user, profile, err := h.authService.Register(req.Username, req.Password, req.FirstName, req.LastName)
	if err != nil {
		// Specific errors
		if err.Error() == "username already exists" {
			log.Warn().Str("username", req.Username).Msg("Username already exists")
			writeErrorResponse(w, http.StatusConflict, "Username already exists")
			return
		}
		if err.Error() == "password must be at least 8 characters long" {
			log.Warn().Msg("Password too short")
			writeErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		log.Error().Err(err).Msg("Failed to register user")
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	// Create reponse
	response := map[string]interface{}{
		"message": "User successfully registered",
		"user": model.UserSummary{
			UserID:   user.ID,
			Username: user.Username,
			Role:     user.Role,
		},
		"profile": profile,
	}

	log.Info().
		Str("username", user.Username).
		Int("user_id", user.ID).
		Msg("User registered successfully")

	writeJSONResponse(w, http.StatusOK, response)
}
