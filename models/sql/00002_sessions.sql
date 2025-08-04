-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID UNIQUE NOT NULL,
    token_hash TEXT UNIQUE NOT NULL,
    CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd