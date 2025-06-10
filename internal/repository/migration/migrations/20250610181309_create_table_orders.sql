-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    number VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    accrual INT NOT NULL,
    uploaded_at TIMESTAMP NOT NULL,
    PRIMARY KEY (number, username)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
