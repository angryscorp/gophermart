-- +goose Up
-- +goose StatementBegin
CREATE TABLE withdrawals (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    order_number VARCHAR(255) NOT NULL,
    withdrawn DECIMAL(15,2) NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL,

    FOREIGN KEY (user_id, order_number) REFERENCES orders(user_id, number) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdrawals;
-- +goose StatementEnd
