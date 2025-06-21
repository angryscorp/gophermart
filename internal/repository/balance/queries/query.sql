-- name: Balance :one
SELECT
    balance,
    withdrawn
FROM
    balances
WHERE
    user_id = @user_id;

-- name: Withdrawals :many
SELECT
    order_number,
    withdrawn,
    processed_at
FROM
    withdrawals
WHERE
    user_id = @user_id;