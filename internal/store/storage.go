package store

import (
	"database/sql"
)

// app.store.Posts.GetById()
type Storage struct {
	Posts PostRepository
	Users UserRepository
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostsStore{db},
		Users: &UsersStore{db},
	}
}
