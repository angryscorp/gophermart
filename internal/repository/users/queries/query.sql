-- name: CreateUser :exec
INSERT INTO users (id, username, password_hash)
VALUES ($1, $2, $3);

-- name: UserByUsername :one
SELECT id, username, password_hash FROM users
WHERE username = $1;

-- name: CheckUsername :one
SELECT EXISTS(
    SELECT 1 FROM users
    WHERE username = $1
) AS user_exists;

