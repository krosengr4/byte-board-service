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
		err := rows.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &post.Content, &post.Author, &post.DatePosted)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		postList = append(postList, post)
	}

	return postList, nil
}

/*
	todo:
		- Add comment
		- Edit comment
		- Delete comment

		- Get all posts
		- Get post by ID
		- Get posts by user ID
		- Add post
		- Update post
		- Delete post

		- Get profile by User ID
		- Create profile
		- Update profile
		- Delete profile

		- Get all users
		- Get user by ID
		- Get user by username
		- Get userID by username
		- Create user
		- Check if user exists
*/
