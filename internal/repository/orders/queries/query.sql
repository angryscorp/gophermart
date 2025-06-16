-- name: GetOrderForUpdate :one
SELECT
    *
FROM
    orders
WHERE
    number == @number
FOR UPDATE;

-- name: CreateOrder :exec
INSERT INTO orders (
    number,
    username,
    status,
    accrual,
    uploaded_at
) VALUES (
    @number,
    @username,
    @status,
    0,
    NOW()
);

-- name: AllOrders :many
SELECT
    *
FROM
    orders
ORDER BY
    uploaded_at DESC;

-- name: UpdateOrder :exec
UPDATE orders
SET
    status = @status,
    accrual = @accrual
WHERE
    number = @number AND
    username = @username;
