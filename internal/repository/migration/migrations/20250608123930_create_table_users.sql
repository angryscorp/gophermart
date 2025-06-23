-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    PRIMARY KEY (username)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
