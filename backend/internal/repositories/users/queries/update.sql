-- name: Update :one
UPDATE users
SET
    name = COALESCE($2, name),
    password_hash = COALESCE($3, password_hash)
WHERE id = $1
RETURNING *;