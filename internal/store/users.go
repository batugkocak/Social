package store

import (
	"context"
	"database/sql"
)

// UserRepository
type UserRepository interface {
	Create(context.Context) error
}

// UserRepository Implementation
type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context) error {
	return nil
}
