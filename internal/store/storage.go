package store

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound       = errors.New("Resource not found.")
	ErrConflict       = errors.New("Resource already exists.")
	QueryTimeDuration = time.Second * 5
)

// app.store.Posts.GetById()
type Storage struct {
	Posts     PostRepository
	Users     UserRepository
	Comments  CommentRepository
	Followers FollowerRepository
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}
