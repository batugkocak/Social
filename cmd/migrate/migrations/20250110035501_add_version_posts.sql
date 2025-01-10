-- +goose Up
-- +goose StatementBegin
ALTER TABLE posts
ADD COLUMN version INT DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE posts
DROP COLUMN version
-- +goose StatementEnd
