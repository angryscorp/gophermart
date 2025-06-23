-- +goose Up
-- +goose StatementBegin
CREATE TABLE balances (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    balance DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    withdrawn DECIMAL(15,2) NOT NULL DEFAULT 0.00,

    CONSTRAINT positive_balance CHECK (balance >= 0),
    CONSTRAINT positive_withdrawn CHECK (withdrawn >= 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS balances;
-- +goose StatementEnd
