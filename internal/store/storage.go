package store

import (
	"context"
	"database/sql"
)

// app.store.Posts.GetById()
type Storage struct {
	Posts interface {
		Create(context.Context) error
	}
	Users interface {
		Create(context.Context) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostsStore{db},
		Users: &UsersStore{db},
	}
}

// Alternative, does the same job.
// type Storage interface {
// 	Posts() PostRepository
// 	Comments() CommentRepository
// 	Users() UserRepository
// 	Followers() FollowerRepository
// 	AuthTokens() AuthTokenRepository
// 	Roles() RoleRepository
// }
