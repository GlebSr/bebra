-- name: Delete :exec
DELETE FROM users
WHERE id = $1;