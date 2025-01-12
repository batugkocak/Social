package store

import (
	"context"
	"database/sql"
)

// Entity
type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

// UserRepository
type UserRepository interface {
	Create(context.Context, *User) error
	GetById(context.Context, int64) (*User, error) 
}

// UserRepository Implementation
type UserStore struct {
	db *sql.DB
}

// Implementation Functions

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, password, email) 
		VALUES($1, $2, $3)
    RETURNING id, created_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetById(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1
	`

	user := &User{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
