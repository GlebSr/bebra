-- name: GetByID :one
SELECT * FROM users
WHERE id = $1;