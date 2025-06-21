-- +goose Up
ALTER TABLE orders ALTER COLUMN accrual DROP NOT NULL;

-- +goose Down
ALTER TABLE orders ALTER COLUMN accrual SET NOT NULL;