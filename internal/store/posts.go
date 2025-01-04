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
type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context) error {

	return nil
}
