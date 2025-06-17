-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders DROP CONSTRAINT orders_pkey;
ALTER TABLE orders DROP COLUMN username;
ALTER TABLE orders ADD COLUMN user_id UUID NOT NULL;
ALTER TABLE orders ADD PRIMARY KEY (number, user_id);
ALTER TABLE orders ADD CONSTRAINT fk_orders_user_id
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_user_id;
ALTER TABLE orders DROP CONSTRAINT orders_pkey;
ALTER TABLE orders DROP COLUMN user_id;
ALTER TABLE orders ADD COLUMN username VARCHAR(255) NOT NULL;
ALTER TABLE orders ADD PRIMARY KEY (number, username);
-- +goose StatementEnd