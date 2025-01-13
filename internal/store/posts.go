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
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comment_count"`
}

// PostRepository
type PostRepository interface {
	Create(context.Context, *Post) error
	GetById(context.Context, int64) (*Post, error)
	DeleteById(context.Context, int64) error
	UpdateById(context.Context, *Post) error
	GetUserFeed(context.Context, int64) ([]*PostWithMetadata, error)
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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags)).Scan(
		&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetById(ctx context.Context, postID int64) (*Post, error) {
	query := `
	SELECT id, user_id, content, title, tags, created_at, updated_at, version
	FROM posts 
	WHERE id = $1` // Change user_id to id

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	fmt.Println(postID)
	var post Post
	err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
		&post.Title,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

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
	SET title = $1, content = $2, tags = $3, updated_at = $4, version = version +1
	WHERE id = $5 AND version = $6
	RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		time.Now(),
		post.ID,
		post.Version,
	).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil

}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64) ([]*PostWithMetadata, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE f.user_id = $1 OR p.user_id = $1
		GROUP BY p.id, u.username
		ORDER BY p.created_at DESC
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []*PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		feed = append(feed, &p)
	}
	return feed, nil
}
