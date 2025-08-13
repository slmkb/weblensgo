-- +goose Up
-- +goose StatementBegin
CREATE TABLE galleries(
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    title TEXT NOT NULL,
    hash TEXT NOT NULL,
    UNIQUE(user_id, title),
    UNIQUE(user_id, hash)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE galleries;
-- +goose StatementEnd