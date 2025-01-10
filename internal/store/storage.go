package store

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound       = errors.New("Resource not found.")
	QueryTimeDuration = time.Second * 5
)

// app.store.Posts.GetById()
type Storage struct {
	Posts    PostRepository
	Users    UserRepository
	Comments CommentRepository
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
