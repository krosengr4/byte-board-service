package repository

import (
	"byte-board/internal/appconfig"
	"byte-board/internal/model"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type DB struct {
	*sql.DB
}

// Create new database connection
func New(cfg *appconfig.Config) (*DB, error) {
	// Get the database URL
	databaseURL, err := cfg.GetDatabaseURL()
	if err != nil {
		return nil, fmt.Errorf("could not get the database url: %w", err)
	}

	// Open connection to database
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("could not establish connection with database: %w", err)
	}

	// Ping database (verify conn to db is still alive)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msg("Database successfully connected!")
	return &DB{DB: db}, nil
}

// #region Comments

// GET api/comments - Get all comments in the db
func (db *DB) GetAllComments() ([]model.Comment, error) {
	query := "SELECT * FROM comments"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var commentsList []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(&comment.CommentId, &comment.UserId, &comment.PostId, &comment.Content, &comment.Author, &comment.DatePosted)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comments: %w", err)
		}

		commentsList = append(commentsList, comment)
	}

	return commentsList, nil
}

// GET api/comment/{commentId} - Get comment by ID
func (db *DB) GetCommentById(commentId int) (*model.Comment, error) {
	query := "SELECT * FROM comments WHERE comment_id = $1"

	var comment model.Comment
	err := db.QueryRow(query, commentId).Scan(&comment.CommentId, &comment.UserId, &comment.PostId, &comment.Content, &comment.Author, &comment.DatePosted)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("comment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}

	return &comment, nil
}

// GET api/post/{postId}/comments - Get all comments on a post
func (db *DB) GetCommentsByPost(postId int) ([]model.Comment, error) {
	query := "SELECT * FROM comments WHERE post_id = $1"

	rows, err := db.Query(query, postId)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments on post: %w", err)
	}
	defer rows.Close()

	var commentList []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(&comment.CommentId, &comment.UserId, &comment.PostId, &comment.Content, &comment.Author, &comment.DatePosted)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comments on post")
		}

		commentList = append(commentList, comment)
	}

	if len(commentList) == 0 {
		return nil, fmt.Errorf("no comments were found")
	}
	return commentList, nil
}

// #endregion

// #region Posts

// GET api/posts - Get all posts in the DB
func (db *DB) GetAllPosts() ([]model.Post, error) {
	query := "SELECT * FROM posts"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query rows: %w", err)
	}
	defer rows.Close()

	var postList []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &post.Author, &post.DatePosted)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		postList = append(postList, post)
	}

	return postList, nil
}

// GET api/posts/{postId} - Get post by post ID
func (db *DB) GetPostById(postId int) (*model.Post, error) {
	query := "SELECT * FROM posts WHERE post_id = $1"

	var post model.Post
	err := db.QueryRow(query, postId).Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &post.Author, &post.DatePosted)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query post with that id: %w", err)
	}

	return &post, nil
}

// GET api/posts/user/{userId} - Get all posts made by a user
func (db *DB) GetPostsByUserId(userId int) ([]model.Post, error) {
	query := "SELECT * FROM posts WHERE user_id = $1"

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query rows: %w", err)
	}

	var postList []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &post.Author, &post.DatePosted)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		postList = append(postList, post)
	}

	if len(postList) == 0 {
		return nil, fmt.Errorf("users posts not found")
	}
	return postList, nil
}

// #endregion

// #region Profiles

// GET api/profiles - Get all profiles
func (db *DB) GetAllProfiles() ([]model.Profile, error) {
	query := "SELECT * FROM profiles"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query profiles: %w", err)
	}

	var profileList []model.Profile
	for rows.Next() {
		var profile model.Profile
		err := rows.Scan(&profile.UserId, &profile.FirstName, &profile.LastName, &profile.Email, &profile.GithubLink, &profile.City, &profile.State, &profile.DateRegistered)
		if err != nil {
			return nil, fmt.Errorf("failed to scan profiles: %w", err)
		}

		profileList = append(profileList, profile)
	}

	return profileList, nil
}

// GET api/profiles/{userId} - Get profile by User ID
func (db *DB) GetProfileByUserId(userId int) (*model.Profile, error) {
	query := "SELECT * FROM profiles WHERE user_id = $1"

	var profile model.Profile
	err := db.QueryRow(query, userId).Scan(&profile.UserId, &profile.FirstName, &profile.LastName, &profile.Email, &profile.GithubLink, &profile.City, &profile.State, &profile.DateRegistered)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("profile not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query profiles: %w", err)
	}

	return &profile, err
}

// #endregion

// #region Users

// GET api/users - Get all users
func (db *DB) GetAllUsers() ([]model.User, error) {
	query := "SELECT * FROM users"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users")
	}

	var userList []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("failed to scan users")
		}

		userList = append(userList, user)
	}

	return userList, nil
}

// GET api/users/{userId} - Get user by user ID
func (db *DB) GetUserByID(userId int) (*model.User, error) {
	query := "SELECT * FROM users WHERE user_id = $1"

	var user model.User
	err := db.QueryRow(query, userId).Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Role)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query or scan rows: %w", err)
	}

	return &user, nil
}

// GET api/users/username/{username} - Get user by username
func (db *DB) GetUserByUsername(username string) (*model.User, error) {
	query := "SELECT * FROM users WHERE username = $1"

	var user model.User
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Role)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("username not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query or scan rows: %w", err)
	}

	return &user, nil
}

// Create new user
func (db *DB) CreateUser(user *model.User) error {
	query := `
		INSERT INTO users (username, hashed_password, role)
		VALUES ($1, $2, $3)
		RETURNING user_id
	`

	err := db.QueryRow(query, user.Username, user.HashedPassword, user.Role).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Update user
func (db *DB) UpdateUser(user *model.User) error {
	query := `
		UPDATE users
		SET username = $1,
		hashed_password = $2,
		role = $3
		WHERE user_id = $4
	`

	result, err := db.Exec(query, user.Username, user.HashedPassword, user.Role, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete user
func (db *DB) DeleteUser(userId int) error {
	query := "DELETE FROM users WHERE user_id = $1"

	result, err := db.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Check if username already exists
func (db *DB) UserExists(username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"

	var exists bool
	err := db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return exists, nil
}

// #endregion

/*
	todo:
		- Add comment
		- Edit comment
		- Delete comment

		- Add post
		- Update post
		- Delete post

		- Get profile by User ID
		- Create profile
		- Update profile
		- Delete profile

		- Get user by username
		- Get userID by username
		- Create user
		- Check if user exists
*/
