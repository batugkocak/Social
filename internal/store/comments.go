package store

import (
	"context"
	"database/sql"
	"time"
)

// Entity
type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `json:"user"`
}

// Repository
type CommentRepository interface {
	GetByPostID(ctx context.Context, postID int64) (*[]Comment, error)
}

// Repository Implementation
type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) (*[]Comment, error) {
	query := `
	SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id
	FROM comments c
	JOIN users on users.id = c.user_id
	WHERE c.post_id = $1
	ORDER BY c.created_at DESC;`

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username, &c.User.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return &comments, nil
}
