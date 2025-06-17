-- name: CreateUser :exec
INSERT INTO users (id, username, password_hash)
VALUES ($1, $2, $3);

-- name: CheckUser :one
SELECT id FROM users
WHERE username = $1 AND password_hash = $2;

-- name: CheckUsername :one
SELECT EXISTS(
    SELECT 1 FROM users
    WHERE username = $1
) AS user_exists;

