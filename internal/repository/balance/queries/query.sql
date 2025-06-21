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

-- name: UpdateBalance :exec
UPDATE balances
SET balance = @balance,
    withdrawn = @withdrawn
WHERE
    user_id = @user_id;

-- name: CheckBalanceForUpdate :one
SELECT
    balance,
    withdrawn
FROM
    balances
WHERE
    user_id = @user_id
FOR UPDATE;

-- name: AddWithdrawal :exec
INSERT INTO withdrawals (id, user_id, order_number, withdrawn, processed_at)
VALUES (@id, @user_id, @order_number, @withdrawn, NOW());