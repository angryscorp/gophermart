-- name: CreateUser :exec
INSERT INTO users (username, password_hash)
VALUES ($1, $2);

-- name: CheckUser :one
SELECT 1 FROM users
WHERE username = $1 AND password_hash = $2;