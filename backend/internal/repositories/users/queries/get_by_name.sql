-- name: GetByName :one
SELECT *
FROM users
WHERE name = $1;