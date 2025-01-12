package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

// Entity
type Follower struct {
	UserID     int64     `json:"user_id"`
	FollowerID int64     `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// Repository
type FollowerRepository interface {
	Follow(ctx context.Context, followedID, userID int64) error
	Unfollow(ctx context.Context, followedID, userID int64) error
}

// Implementation
type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerID, userID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return ErrConflict
			}
		}
		return err
	}
	return err
}

func (s *FollowerStore) Unfollow(ctx context.Context, followerID, userID int64) error {
	query := `
		DELETE FROM followers 
		WHERE user_id = $1 AND follower_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	return err
}
