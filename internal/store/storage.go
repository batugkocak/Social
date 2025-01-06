package store

import (
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("Resource not found.")
)

// app.store.Posts.GetById()
type Storage struct {
	Posts PostRepository
	Users UserRepository
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db},
		Users: &UserStore{db},
	}
}
