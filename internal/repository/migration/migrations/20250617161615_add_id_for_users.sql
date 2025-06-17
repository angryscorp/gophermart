-- +goose Up
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users ADD COLUMN id UUID PRIMARY KEY;
CREATE UNIQUE INDEX users_username_idx ON users (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS users_username_idx;
ALTER TABLE users DROP COLUMN IF EXISTS id;
ALTER TABLE users ADD CONSTRAINT users_pkey PRIMARY KEY (username);
-- +goose StatementEnd