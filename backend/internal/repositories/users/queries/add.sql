-- name: Add :one
INSERT INTO users (
    id, name, password_hash
) VALUES (
    $1, $2, $3
)
RETURNING *;