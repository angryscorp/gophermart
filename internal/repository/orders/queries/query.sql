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
    user_id,
    status,
    accrual,
    uploaded_at
) VALUES (
    @number,
    @user_id,
    @status,
    0,
    NOW()
);

-- name: AllOrders :many
SELECT
    *
FROM
    orders
WHERE
    user_id = @user_id
ORDER BY
    uploaded_at DESC;

-- name: UpdateOrder :exec
UPDATE orders
SET
    status = @status,
    accrual = @accrual
WHERE
    number = @number AND
    user_id = @user_id;

-- name: IncreaseBalance :exec
UPDATE balances
SET
    balance = balance + orders.accrual
FROM
    orders
WHERE
    balances.user_id = orders.user_id
  AND
    orders.number = @number;

