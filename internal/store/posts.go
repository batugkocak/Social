package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// Entity
type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserId    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

// PostRepository
type PostRepository interface {
	Create(context.Context, *Post) error
	GetById(context.Context, int64) (*Post, error)
	DeleteById(context.Context, int64) error
	UpdateById(ctx context.Context, post *Post) error
}

// PostRepository Implementation
type PostStore struct {
	db *sql.DB
}

// Implementation Functions

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags)
	VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserId, pq.Array(post.Tags)).Scan(
		&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetById(ctx context.Context, postID int64) (*Post, error) {
	query := `
	SELECT id, user_id, content, title, tags, created_at, updated_at
	FROM posts 
	WHERE id = $1` // Change user_id to id

	fmt.Println(postID)
	var post Post
	err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.UserId,
		&post.Content,
		&post.Title,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *PostStore) DeleteById(ctx context.Context, postID int64) error {
	query := `
	DELETE FROM posts 
	WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostStore) UpdateById(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts
	SET title = $1, content = $2, tags = $3, updated_at = $4
	WHERE id = $5
	`
	result, err := s.db.ExecContext(ctx, query, post.Title, post.Content, pq.Array(post.Tags), time.Now(), post.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrNotFound
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil

}
