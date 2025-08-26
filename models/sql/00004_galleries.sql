-- +goose Up
-- +goose StatementBegin
CREATE TABLE galleries(
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    title TEXT NOT NULL,
    hash TEXT NOT NULL UNIQUE,
    UNIQUE(user_id, title)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE galleries;
-- +goose StatementEnd