-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD CONSTRAINT orders_number_unique UNIQUE (number);
ALTER TABLE withdrawals DROP CONSTRAINT IF EXISTS withdrawals_user_id_order_number_fkey;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE withdrawals DROP CONSTRAINT IF EXISTS withdrawals_order_number_fkey;
ALTER TABLE withdrawals ADD CONSTRAINT withdrawals_user_id_order_number_fkey
    FOREIGN KEY (user_id, order_number) REFERENCES orders(user_id, number) ON DELETE CASCADE;
-- +goose StatementEnd