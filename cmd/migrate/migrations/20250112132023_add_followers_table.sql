-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS followers(
  user_id bigint NOT NULL,
  follower_id bigint NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

  PRIMARY KEY(user_id, follower_id), -- Composite Key
  FOREIGN KEY(user_id) REFERENCES users (id) ON DELETE CASCADE, -- DELETE CASCADE is not good for soft delete ofc.
  FOREIGN KEY(follower_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS followers
-- +goose StatementEnd
