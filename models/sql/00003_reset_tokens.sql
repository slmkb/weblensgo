-- +goose Up
-- +goose StatementBegin
CREATE TABLE reset_tokens (
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE reset_tokens;
-- +goose StatementEnd