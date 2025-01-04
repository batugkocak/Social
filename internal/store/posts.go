package store

import (
	"context"
	"database/sql"
)

// PostRepository
type PostRepository interface {
	Create(context.Context) error
}

// PostRepository Implementation
type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context) error {
	return nil
}
