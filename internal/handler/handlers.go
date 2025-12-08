package handler

import (
	"byte-board/internal/appconfig"
	"byte-board/internal/middleware"
	"byte-board/internal/model"
	"byte-board/internal/repository"
	"byte-board/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	db          *repository.DB
	config      *appconfig.Config
	authService *service.AuthService
}

// Create a new instance of a handler
func New(db *repository.DB, cfg *appconfig.Config, authService *service.AuthService) *Handler {
	return &Handler{
		db:          db,
		config:      cfg,
		authService: authService,
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

// #region Comment handlers

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
		log.Warn().Str("id", idStr).Msg("Invalid ID format")
		writeErrorResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Get comment by id from the database
	comment, err := h.db.GetCommentById(id)
	if err != nil {
		if err.Error() == "comment not found" {
			log.Warn().Int("ID", id).Msg("Comment with that ID not found")
			writeErrorResponse(w, http.StatusNotFound, "Comment not found")
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
		log.Error().Err(err).Msg("Failed to get all comments on the post")
		writeErrorResponse(w, http.StatusInternalServerError, "failed to get comments on post")
		return
	}

	log.Info().Int("count", len(comments)).Msg("Successfully retrieved comments on post")
	writeJSONResponse(w, http.StatusOK, comments)

}

// #endregion

// #region Post handlers

// Handler to get all posts
func (h *Handler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET /posts - Getting all posts")

	posts, err := h.db.GetAllPosts()
	if err != nil {
		log.Error().Err(err).Msg("Error getting all posts")
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get all posts")
		return
	}

	log.Info().Int("count", len(posts)).Msg("Successfully retrieved all posts")
	writeJSONResponse(w, http.StatusOK, posts)
}

// Handler to get post by ID
func (h *Handler) GetPostById(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET /posts/{postId} - Getting a post by post ID")

	vars := mux.Vars(r)
	idStr := vars["postId"]

	// Convert the ID from string to an int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().Str("ID", idStr).Msg("Invalid post ID format")
		writeErrorResponse(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	post, err := h.db.GetPostById(id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get post by ID")
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get post by ID")
		return
	}

	log.Info().Int("Post ID", id).Msg("Successfully retrieved post by ID")
	writeJSONResponse(w, http.StatusOK, post)
}

// Handler to get all posts by UserID
func (h *Handler) GetPostsByUserId(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET /posts/user/{userId} - Getting all posts by user ID")

	vars := mux.Vars(r)
	idStr := vars["userId"]

	// Convert string ID into an int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().Str("ID", idStr).Msg("Invalid user ID format")
		writeErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	posts, err := h.db.GetPostsByUserId(id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get posts from that user")
		writeErrorResponse(w, http.StatusInternalServerError, "Failure to get posts with that user ID")
		return
	}

	log.Info().Int("Count", len(posts)).Msg("Successfully retrieved posts from user ID")
	writeJSONResponse(w, http.StatusOK, posts)
}

// POST /api/protected/posts - Create new post
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("POST /api/posts - Creating new post")

	// Get authenticated user from JWT mware context
	username := middleware.GetUsername(r)
	if username == "" {
		log.Warn().Msg("No username in context")
		writeErrorResponse(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Get user from db
	user, err := h.db.GetUserByUsername(username)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user")
		writeErrorResponse(w, http.StatusInternalServerError, "failed to get user info")
		return
	}

	// Parse body request
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Title == "" || req.Content == "" {
		log.Warn().Msg("Missing required fields")
		writeErrorResponse(w, http.StatusBadRequest, "Title and content are required")
		return
	}

	// Create post object
	post := &model.Post{
		UserId:     user.ID,
		Title:      req.Title,
		Content:    req.Content,
		Author:     user.Username,
		DatePosted: time.Now(),
	}

	// Call db to create post
	if err := h.db.CreatePost(post); err != nil {
		log.Error().Err(err).Msg("failed to create post")
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	log.Info().Str("title", post.Title).Msg("Post created successfully")
	writeJSONResponse(w, http.StatusCreated, post)
}

// #endregion

// #region Profile handlers

// Handler to get all profiles
func (h *Handler) GetAllProfiles(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET /profiles - Getting all profiles")

	profiles, err := h.db.GetAllProfiles()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all profiles")
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get profiles")
		return
	}

	log.Info().Int("Count", len(profiles)).Msg("Successfully retrieved all profiles")
	writeJSONResponse(w, http.StatusOK, profiles)
}

// Handler to get profile by User ID
func (h *Handler) GetProfileByUserId(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET /profiles/{userId} - Getting profile by user ID")

	// Get userID
	vars := mux.Vars(r)
	idStr := vars["userId"]

	// Convert string user ID to an int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().Str("User ID", idStr).Msg("Invalid user ID format")
		writeErrorResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	profile, err := h.db.GetProfileByUserId(id)
	if err != nil {
		if err.Error() == "profile not found" {
			log.Warn().Int("ID", id).Msg("Profile not found")
			writeErrorResponse(w, http.StatusNotFound, "Profile not found")
			return
		}
		log.Error().Err(err).Msg("Error getting profile")
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get profile")
		return
	}

	log.Info().Int("ID", id).Msg("Successfully retrieved profile")
	writeJSONResponse(w, http.StatusOK, profile)
}

// #endregion

// #region Handler for Users

// Handler to get all Users
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET /users - Getting all users")

	users, err := h.db.GetAllUsers()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all users")
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get all users")
		return
	}

	log.Info().Msg("Successfully retrieved all users")
	writeJSONResponse(w, http.StatusOK, users)
}

// Handler to get User by User ID
func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET /users/{userId} - Getting user by user ID")

	// Get ID
	vars := mux.Vars(r)
	idStr := vars["userId"]

	// Convert int UserID to a string
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().Str("ID", idStr).Msg("Invalid user ID format")
		writeErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.db.GetUserByID(id)
	if err != nil {
		if err.Error() == "user not found" {
			log.Warn().Int("ID", id).Msg("No user with that ID found")
			writeErrorResponse(w, http.StatusNotFound, "User not found")
			return
		}
		log.Error().Err(err).Msg("Failed to get user with that ID")
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	log.Info().Int("ID", id).Msg("Successfully retrieved user")
	writeJSONResponse(w, http.StatusOK, user)
}

// Handler to get User by Username
func (h *Handler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET /users/username/{username} - Getting user by username")

	// Get username
	vars := mux.Vars(r)
	username := vars["username"]

	user, err := h.db.GetUserByUsername(username)
	if err != nil {
		if err.Error() == "username not found" {
			log.Warn().Str("username", username).Msg("No user with that username found")
			writeErrorResponse(w, http.StatusNotFound, "Username not found")
			return
		}
		log.Error().Err(err).Msg("Failed to get user with that username")
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	log.Info().Str("Username", username).Msg("Successfully retrieved user")
	writeJSONResponse(w, http.StatusOK, user)
}

// #endregion
